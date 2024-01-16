package main

import (
	"github.com/gin-gonic/gin"

	"github.com/dariomatias-dev/go_auth/api/routes"
	"github.com/dariomatias-dev/go_auth/initialize"
)

func init() {
	initialize.Load()
}

func main() {
	app := gin.Default()

	routes.AppRoutes(app)

	app.Run("localhost:3001")
}
