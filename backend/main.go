package main

import (
	"os"
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

	select {} //! REMOVE THIS LINE
}
