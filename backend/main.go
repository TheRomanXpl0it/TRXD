package main

import (
	"context"
	"flag"
	"os"
	"strings"
	"trxd/api"
	"trxd/api/routes/user_register"
	"trxd/db"
	"trxd/db/sqlc"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/joho/godotenv"
	"github.com/tde-nico/log"
)

func Flags() {
	var (
		help                bool
		h                   bool
		user                string
		toggleRegisterAllow bool
	)
	flag.BoolVar(&help, "help", false, "Show help")
	flag.BoolVar(&h, "h", false, "Show help")
	flag.BoolVar(&toggleRegisterAllow, "t", false, "Toggle the allow-register config")
	flag.StringVar(&user, "r", "", "Register a new admin user with 'username:email:password'")
	flag.Parse()

	if help || h {
		flag.Usage()
		os.Exit(0)
	}

	if toggleRegisterAllow {
		err := db.UpdateConfig(context.Background(), "allow-register", "true")
		if err != nil {
			log.Fatal("Error updating allow-register config", "err", err)
		}
		os.Exit(0)
	}

	if user != "" {
		parts := strings.SplitN(user, ":", 3)
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

		user, err := user_register.RegisterUser(context.Background(), name, email, password, sqlc.UserRoleAdmin)
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
	if _, err := os.Stat("DEV"); !os.IsNotExist(err) {
		log.SetLogLevel("debug")
		log.SetReportCaller(false)
	}

	godotenv.Load()

	info, err := utils.GetDBInfoFromEnv()
	if err != nil {
		log.Fatal("Error getting database info from env", "err", err)
	}

	err = db.ConnectDB(info)
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
