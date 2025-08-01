package main

import (
	"context"
	"flag"
	"os"
	"strings"
	"trxd/api"
	"trxd/db"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/joho/godotenv"
	"github.com/tde-nico/log"
)

func Flags() {
	var (
		help     bool
		h        bool
		register string
	)
	flag.BoolVar(&help, "help", false, "Show help")
	flag.BoolVar(&h, "h", false, "Show help")
	flag.StringVar(&register, "r", "", "Register a new admin user with 'username:email:password'")
	flag.Parse()

	if help || h {
		flag.Usage()
		os.Exit(0)
	}

	if register != "" {
		parts := strings.SplitN(register, ":", 3)
		var name, email, password string
		if len(parts) == 2 {
			var err error
			password, err = utils.GenerateRandPass()
			if err != nil {
				log.Fatal("Error generating random password", "err", err)
			}
			log.Warn("No password provided, using generated password:", "password", password)
			name, email = parts[0], parts[1]
		} else if len(parts) == 3 {
			name, email, password = parts[0], parts[1], parts[2]
		} else {
			log.Fatal("Invalid format for registration. Use 'username:email:password'")
		}

		if name == "" || email == "" || password == "" {
			log.Fatal("Username, email, and password must not be empty")
		}
		if len(password) < consts.MinPasswordLength {
			log.Fatal(consts.ShortPassword)
		}
		if len(name) > consts.MaxNameLength {
			log.Fatal(consts.LongName)
		}
		if len(email) > consts.MaxEmailLength {
			log.Fatal(consts.LongEmail)
		}
		if len(password) > consts.MaxPasswordLength {
			log.Fatal(consts.LongPassword)
		}

		user, err := db.RegisterUser(context.Background(), name, email, password, db.UserRoleAdmin)
		if err != nil {
			log.Fatal("Error registering admin user", "err", err)
		}
		if user == nil {
			log.Fatal("Failed to register admin user: user already exists")
		}
		log.Info("Admin user registered successfully")
		os.Exit(0)
	}
}

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

	Flags()

	log.Info("Starting TRXd server")
	defer log.Info("Stopping TRXd server")

	app := api.SetupApp()
	err = app.Listen(":1337")
	if err != nil {
		log.Fatal("Error starting server", "err", err)
	}
}
