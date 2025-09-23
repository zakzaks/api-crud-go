package main

import (
	"api-crud-go/internal/database"
	"api-crud-go/internal/env"
	"database/sql"
	"log"
	_ "api-crud-go/docs"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
)

// @title Go Gin REST API
// @version 1.0
// @description This is a sample server for a CRUD API built with Go and Gin framework.
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @desciption Enter your beare token in the format  **Bearer &lt;token&gt;**
// @termsOfService http://swagger.io/terms/

type application struct {
	port int
	jwtSecret string
	models database.Models
}

func main(){
	db, err := sql.Open("sqlite3","./data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	models := database.NewModels(db)
	app := &application{
		port:     env.GetEnvInt("PORT", 8080),
		jwtSecret: env.GetEnv("JWT_SECRET", "mysecret"),
		models:   models,
	}

	if err := app.serve(); err != nil {
		log.Fatal(err)
	}
}