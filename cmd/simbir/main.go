package main

import (
	"github.com/joho/godotenv"
	"goSimbir/internal/models"
	"goSimbir/internal/server"
	"goSimbir/pkg/db"
	"log"
	"os"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	// INFO: init Config
	database, _ := os.LookupEnv("POSTGRES_DATABASE")
	user, _ := os.LookupEnv("POSTGRES_USER")
	password, _ := os.LookupEnv("POSTGRES_PASSWORD")

	dbCfg := db.Config{
		User:     user,
		Password: password,
		DBName:   database,
	}
	srvCfg := server.Config{
		Host: "",
		Port: "8080",
	}

	// INFO: init DB
	db, err := db.InitDB(dbCfg)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	// INFO: init models (Domain model)
	models.SetDB(db)

	// INFO: init server
	if err := server.InitServer(srvCfg); err != nil {
		log.Panic(err)
	}
}
