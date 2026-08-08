package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"kopibang/api/route"
	"kopibang/bootstrap"
)

// @title Kopibang Coffee API
// @version 1.0
// @description API Documentation for Kopibang Coffee Point of Sales & Customer App.
// @termsOfService http://swagger.io/terms/

// @contact.name Sefza
// @contact.email admin@sefza.com

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer " followed by a space and your JWT token.

func main() {
	app := bootstrap.App()
	env := app.Env
	defer app.CloseDBConnection()

	timeout := time.Duration(env.ContextTimeout) * time.Second
	ginEngine := gin.Default()

	ginEngine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	route.Setup(env, timeout, &app, ginEngine)
	ginEngine.Run(env.ServerAddress)
}