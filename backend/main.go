package main

import (
	"os"
	"strings"
	"trxd/api"
	"trxd/db"

	"github.com/joho/godotenv"
	"github.com/tde-nico/log"
)

func main() {
	godotenv.Load()

	err := db.ConnectDB(
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)
	if err != nil {
		log.Fatal("Error connecting to database", "err", err)
	}
	defer db.CloseDB()

	err = db.ExecSQLFile("sql/schema.sql")
	if err != nil {
		log.Fatal("Error executing schema SQL", "err", err)
	}
	files, err := os.ReadDir("sql/triggers")
	if err != nil {
		log.Fatal("Error reading triggers directory", "err", err)
	}
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}
		err = db.ExecSQLFile("sql/triggers/" + file.Name())
		if err != nil {
			log.Fatal("Error executing trigger SQL", "file", file.Name(), "err", err)
		}
	}

	log.Info("Starting TRXd server")
	defer log.Info("Stopping TRXd server")

	app := api.SetupApp()
	err = app.Listen(":1337")
	if err != nil {
		log.Fatal("Error starting server", "err", err)
	}
}
