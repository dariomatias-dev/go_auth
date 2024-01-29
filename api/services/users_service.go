package services

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/mail"
	"net/smtp"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	db "github.com/dariomatias-dev/go_auth/api/db/sqlc"
	"github.com/dariomatias-dev/go_auth/api/models"
	"github.com/dariomatias-dev/go_auth/api/utils"
)

type verificationEmailResponse struct {
	Message string `json:"message"`
	Error   error  `json:"error"`
}

type UsersService struct {
	DbQueries *db.Queries
}

func (us UsersService) Create(
	ctx *gin.Context,
	createUserBody models.CreateUserModel,
	userRoles []string,
) *uuid.UUID {
	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(createUserBody.Password),
		10,
	)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			err.Error(),
		)
		return nil
	}

	// Create User Table
	createUserParams := db.CreateUserParams{
		Name:     createUserBody.Name,
		Age:      createUserBody.Age,
		Email:    createUserBody.Email,
		Password: string(encryptedPassword),
		Roles: userRoles,
	}

	userID, err := us.DbQueries.CreateUser(ctx, createUserParams)
	if err != nil {
		panic(err)
	}

	// Create Tokens Table
	err = us.DbQueries.CreateTokens(
		ctx,
		userID,
	)
	if err != nil {
		panic(err)
	}

	// Create Login Attempt Table
	err = us.DbQueries.CreateLoginAttempts(
		ctx,
		userID,
	)
	if err != nil {
		panic(err)
	}

	return &userID
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

func (us UsersService) FindOneByEmail(
	ctx *gin.Context,
	email string,
) *db.Users {
	user, err := us.DbQueries.GetUserByEmail(ctx, email)

	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
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
	updateUserBody models.UpdateModel,
) {
	getValue := utils.GetValue{}

	password := updateUserBody.Password

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
		Name:     getValue.String(updateUserBody.Name),
		Age:      getValue.Int32(updateUserBody.Age),
		Email:    getValue.String(updateUserBody.Email),
		Password: getValue.String(password),
	}

	err := us.DbQueries.UpdateUser(ctx, updateUserParams)
	if err != nil {
		panic(err)
	}
}

func (us UsersService) Delete(
	ctx *gin.Context,
	ID uuid.UUID,
) {
	err := us.DbQueries.DeleteTokens(ctx, ID)
	if err != nil {
		panic(err)
	}

	err = us.DbQueries.DeleteLoginAttempt(ctx, ID)
	if err != nil {
		panic(err)
	}

	err = us.DbQueries.DeleteUser(ctx, ID)
	if err != nil {
		panic(err)
	}
}

func (us UsersService) SendVerificationEmail(
	verificationCode string,
	userName string,
	userEmail string,
) verificationEmailResponse {
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
		return verificationEmailResponse{
			Message: "Error sending e-mail",
			Error:   err,
		}
	}

	return verificationEmailResponse{
		Message: "E-mail successfully sent",
		Error:   nil,
	}
}

func (us UsersService) CreateEmailValidation(
	ctx *gin.Context,
	emailValidation models.EmailValidationModel,
) {
	createEmailValidationParams := db.CreateEmailValidationParams{
		UserID:           emailValidation.UserID,
		VerificationCode: emailValidation.VerificationCode,
		ExpirationTime:   int32(emailValidation.ExpirationTime),
	}

	err := us.DbQueries.CreateEmailValidation(ctx, createEmailValidationParams)
	if err != nil {
		panic(err)
	}
}
