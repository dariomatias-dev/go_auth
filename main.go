package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	db "github.com/dariomatias-dev/go_auth/api/db/sqlc"
	"github.com/dariomatias-dev/go_auth/api/routes"
	"github.com/dariomatias-dev/go_auth/initialize"
)

var dbQueries *db.Queries

func init() {
	initialize.Load()

	dbName := os.Getenv("DATABASE_DRIVER")
	dbURL := fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s?sslmode=disable",
		dbName,
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	dbcon, err := sql.Open(dbName, dbURL)

	if err != nil {
		log.Fatal(err)
	}

	if err := dbcon.Ping(); err != nil {
		log.Fatal(err)
	}

	dbQueries = db.New(dbcon)
}

func main() {
	app := gin.Default()

	routes.AppRoutes(
		app,
		dbQueries,
	)

	app.Run("localhost:3001")
}
