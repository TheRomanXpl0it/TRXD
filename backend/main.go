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

	//! TESTS

	// user, err := db.RegisterUser("test", "test@test.test", "test")
	// if err != nil {
	// 	log.Fatal("Error registering user", "err", err)
	// }
	// if user == nil {
	// 	log.Warn("User already exists")
	// 	return
	// }
	// log.Info("User registered", "user", user)

	// team, err := db.RegisterTeam("test", "test", user.ID)
	// if err != nil {
	// 	log.Fatal("Error registering team", "err", err)
	// }
	// if team == nil {
	// 	log.Warn("Team already exists")
	// 	return
	// }
	// log.Info("Team registered", "team", team)

	// user, err := db.LoginUser("test@test.test", "test")
	// if err != nil {
	// 	log.Fatal("Error logging in user", "err", err)
	// }
	// if user == nil {
	// 	log.Warn("User not found or invalid password")
	// 	return
	// }
	// log.Info("User logged in", "user", user)

	// for _, flag := range []string{"flag1", "flag2", "flag3", "aa", "aa-123", "flag", "test"} {
	// 	valid, err := db.SubmitFlag(user.ID, 54, flag)
	// 	if err != nil {
	// 		log.Fatal("Error submitting flag", "err", err)
	// 	}
	// 	log.Info("Flag submission result", "valid", valid)
	// }

	select {} //! REMOVE THIS LINE
}
