package main

import (
	"os"
	"trxd/api"
	"trxd/db"

	"github.com/joho/godotenv"
	"github.com/tde-nico/log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", "err", err)
	}

	err = db.ConnectDB(
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)
	if err != nil {
		log.Fatal("Error connecting to database", "err", err)
	}
	defer db.CloseDB()
	log.Info("Database connection established")

	err = db.ExecSQLFile("sql/schema.sql")
	if err != nil {
		log.Fatal("Error executing schema SQL", "err", err)
	}
	err = db.ExecSQLFile("sql/triggers.sql")
	if err != nil {
		log.Fatal("Error executing triggers SQL", "err", err)
	}

	app := api.SetupApp()
	err = app.Listen(":1337")
	if err != nil {
		log.Fatal("Error starting server", "err", err)
	}
}
