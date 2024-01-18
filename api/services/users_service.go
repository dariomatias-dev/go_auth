package services

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/mail"
	"net/smtp"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	db "github.com/dariomatias-dev/go_auth/api/db/sqlc"
	usertype "github.com/dariomatias-dev/go_auth/api/enums/user_type"
	"github.com/dariomatias-dev/go_auth/api/models"
	"github.com/dariomatias-dev/go_auth/api/utils"
)

type UsersService struct {
	DbQueries *db.Queries
}

func (us UsersService) Create(
	ctx *gin.Context,
	createUser models.CreateUserModel,
) *db.CreateUserRow {
	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(createUser.Password),
		10,
	)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			err.Error(),
		)
		return nil
	}

	// Create User
	createUserParams := db.CreateUserParams{
		Name:     createUser.Name,
		Age:      createUser.Age,
		Email:    createUser.Email,
		Password: string(encryptedPassword),
		Roles: []string{
			usertype.User,
		},
	}

	createdUser, err := us.DbQueries.CreateUser(ctx, createUserParams)

	if err != nil {
		panic(err)
	}

	// Create Tokens
	createTokensParams := db.CreateTokensParams{
		UserID:       createdUser.ID,
		AccessToken:  "",
		RefreshToken: "",
	}

	_, err = us.DbQueries.CreateTokens(ctx, createTokensParams)

	if err != nil {
		panic(err)
	}

	return &createdUser
}

func (us UsersService) FindOne(
	ctx *gin.Context,
	ID uuid.UUID,
) *db.Users {
	user, err := us.DbQueries.GetUser(ctx, ID)
	if err != nil {
		panic(err)
	}

	return &user
}

func (us UsersService) FindAll(
	ctx *gin.Context,
) *[]db.GetUsersRow {
	users, err := us.DbQueries.GetUsers(ctx)
	if err != nil {
		panic(err)
	}

	return &users
}

func (us UsersService) Update(
	ctx *gin.Context,
	ID uuid.UUID,
	updateUser models.UpdateModel,
) *db.UpdateUserRow {
	getValue := utils.GetValue{}

	password := updateUser.Password

	if password != nil {
		value, err := bcrypt.GenerateFromPassword(
			[]byte(*password),
			10,
		)
		if err != nil {
			panic(err)
		}

		encryptedPassword := string(value)
		password = &encryptedPassword
	}

	updateUserParams := db.UpdateUserParams{
		ID:       ID,
		Name:     getValue.String(updateUser.Name),
		Age:      getValue.Int32(updateUser.Age),
		Email:    getValue.String(updateUser.Email),
		Password: getValue.String(password),
	}

	updatedUser, err := us.DbQueries.UpdateUser(ctx, updateUserParams)
	if err != nil {
		panic(err)
	}

	return &updatedUser
}

func (us UsersService) Delete(
	ctx *gin.Context,
	ID uuid.UUID,
) *db.DeleteUserRow {
	_, err := us.DbQueries.DeleteTokens(ctx, ID)
	if err != nil {
		panic(err)
	}

	deletedUser, err := us.DbQueries.DeleteUser(ctx, ID)
	if err != nil {
		panic(err)
	}

	return &deletedUser
}

func (us UsersService) SendVerificationEmail(
	userName string,
	userEmail string,
) string {
	verificationCode := ""
	for loop := 0; loop < 6; loop++ {
		verificationCode += fmt.Sprint(rand.Intn(10))
	}

	from := mail.Address{
		Name:    "Go Auth",
		Address: os.Getenv("GO_AUTH_EMAIL"),
	}
	to := mail.Address{
		Name:    userName,
		Address: userEmail,
	}
	subject := "Activation Code"

	emailBodyHeader := `
		<head>
			<style>
				.no-reply {
					font-weight: 600;
					font-size: 14px;
				}

				.message {
					color: #FFFFFF;
					font-size: 16px;
				}

				.verification-code {
					font-size: 26px;
					font-weight: 800;
					letter-spacing: 4px;
				}

				.copyright-message {
					font-size: 14px;
					color: #6F6F6F;
				}
			</style>
		</head>
	`
	body := fmt.Sprintf(
		`
			<html>
				%s
				<body>
					<main>
						<h2 class="title">
							Activate your account
						</h2>
						<p class="no-reply">
							Do not respond to this email
						</p>
						<p class="message">
							Hello %s.
						</p>
						<p class="message">
							To start using our services, activate your account using the following code:
						</p>
						<h1 class="verification-code">
							%s
						</h1>
						<p class="copyright-message">
							Copyright © 2024. All rights reserved.
						</p>
					</main>
				</body>
			</html>
		`,
		emailBodyHeader,
		userName,
		verificationCode,
	)

	headers := map[string]string{
		"From":         from.String(),
		"To":           to.String(),
		"Subject":      subject,
		"MIME-version": "1.0",
		"Content-Type": "text/html; charset=\"UTF-8\";",
	}

	message := ""
	for key, value := range headers {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	message += "\r\n" + body

	auth := smtp.PlainAuth(
		"",
		from.Address,
		os.Getenv("GO_AUTH_EMAIL_APP_PASSWORD"),
		"smtp.gmail.com",
	)

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		from.Address,
		[]string{
			to.Address,
		},
		[]byte(message),
	)

	if err != nil {
		return fmt.Sprint("Error sending email: ", err)
	}

	return "Email successfully sent"
}
