package initialize

import (
	"context"
	"database/sql"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"

	db "github.com/dariomatias-dev/go_auth/api/db/sqlc"
	usertype "github.com/dariomatias-dev/go_auth/api/enums/user_type"
)

func InitializeAdminUser(
	dbQueries *db.Queries,
) {
	ctx := context.Background()

	adminEmail := os.Getenv("ADMIN_EMAIL")

	_, err := dbQueries.GetUserByEmail(
		ctx,
		adminEmail,
	)

	if err == sql.ErrNoRows {
		encryptedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(os.Getenv("ADMIN_PASSWORD")),
			10,
		)
		if err != nil {
			log.Fatal(err)
		}

		createUserParams := db.CreateUserParams{
			Name:     "Administrator",
			Age:      18,
			Email:    adminEmail,
			Password: string(encryptedPassword),
			Roles: []string{
				usertype.Admin,
			},
		}

		dbQueries.CreateUser(
			ctx,
			createUserParams,
		)
	} else if err != nil {
		panic(err)
	}
}
