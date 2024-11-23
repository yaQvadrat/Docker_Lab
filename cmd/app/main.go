package main

import (
	"app/internal/app"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		return
	}
}

func main() {
	app.Migrations()
	app.Run()
}
