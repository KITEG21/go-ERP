package main

import (
	"os"
	_ "user_api/cmd/docs"
	"user_api/internal/app"
	"user_api/internal/database"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load("../.env")
	database.Connect()
	if err := database.RunMigrations(); err != nil {
		panic(err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default_secret"
	}

	a, err := app.NewApp(jwtSecret)
	if err != nil {
		panic(err)
	}
	a.Engine.Run(":8080")
}
