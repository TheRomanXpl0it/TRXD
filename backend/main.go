package main

import (
	"context"
	"flag"
	"os"
	"strings"
	"trxd/api"
	"trxd/api/routes/users_register"
	"trxd/db"
	"trxd/db/sqlc"
	"trxd/instancer"
	"trxd/utils"
	"trxd/utils/consts"
	"trxd/utils/crypto_utils"
	"trxd/validator"

	"github.com/joho/godotenv"
	"github.com/tde-nico/log"
)

func toggleRegister(ctx context.Context) {
	conf, err := db.GetConfig(ctx, "allow-register")
	if err != nil {
		log.Fatal("Error getting allow-register config", "err", err)
	}
	if conf == "" {
		log.Fatal("allow-register config not found")
	}

	var toggle string
	if conf == "false" {
		toggle = "true"
	} else {
		toggle = "false"
	}

	err = db.UpdateConfig(ctx, "allow-register", toggle)
	if err != nil {
		log.Fatal("Error updating allow-register config", "err", err)
	}

	log.Notice("allow-register set to:", "value", toggle)
}

func registerAdmin(ctx context.Context, userInfo string) {
	parts := strings.SplitN(userInfo, ":", 3)
	var name, email, password string
	if len(parts) == 2 {
		var err error
		password, err = crypto_utils.GeneratePassword()
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
	valid, err := validator.Var(nil, name, "user_name")
	if err != nil || !valid {
		log.Fatal(err)
	}
	valid, err = validator.Var(nil, email, "user_email")
	if err != nil || !valid {
		log.Fatal(err)
	}
	valid, err = validator.Var(nil, password, "password")
	if err != nil || !valid {
		log.Fatal(err)
	}

	tx, err := db.BeginTx(ctx)
	if err != nil {
		log.Fatal("Error beginning transaction", "err", err)
	}
	defer tx.Rollback()

	user, err := users_register.RegisterUser(ctx, tx, name, email, password, sqlc.UserRoleAdmin)
	if err != nil {
		log.Fatal("Error registering admin user", "err", err)
	}
	if user == nil {
		log.Fatal("Failed to register admin user: user already exists")
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal("Error committing transaction", "err", err)
	}

	log.Info("Admin user registered successfully")
}

func flushCache(ctx context.Context) {
	err := db.StorageFlush(ctx)
	if err != nil {
		log.Fatal("Error flushing the cache", "err", err)
	}
}

func insertTestData(ctx context.Context) {
	log.Warn("Inserting mock data into the database. This will delete all existing data!")

	_, err := db.ExecSQLFile("sql/tests.sql")
	if err != nil {
		log.Fatal("Error executing SQL file", "err", err)
	}

	err = db.DeleteAll(ctx)
	if err != nil {
		log.Fatal("Error deleting existing data", "err", err)
	}
	err = db.InitConfigs()
	if err != nil {
		log.Fatal("Error initializing configs", "err", err)
	}
	err = db.InsertMockData()
	if err != nil {
		log.Fatal("Error inserting mock data", "err", err)
	}
}

func parseFlags(ctx context.Context) {
	var (
		help               bool
		h                  bool
		user               string
		toggleRegisterFlag bool
		flushCacheFlag     bool
		insertTestDataFlag bool
	)
	flag.BoolVar(&help, "help", false, "Show help")
	flag.BoolVar(&h, "h", false, "Show help")
	flag.BoolVar(&toggleRegisterFlag, "t", false, "Toggle the allow-register config")
	flag.StringVar(&user, "r", "", "Register a new admin user with 'username:email:password'")
	flag.BoolVar(&flushCacheFlag, "f", false, "Flush the system cache")
	flag.BoolVar(&insertTestDataFlag, "test-data-WARNING-DO-NOT-USE-IN-PRODUCTION", false, "Inserts mocks data into the db")
	flag.Parse()

	switch {
	case help || h:
		flag.Usage()
	case toggleRegisterFlag:
		toggleRegister(ctx)
	case user != "":
		registerAdmin(ctx, user)
	case flushCacheFlag:
		flushCache(ctx)
	case insertTestDataFlag:
		insertTestData(ctx)
	default:
		return
	}

	os.Exit(0)
}

func main() {
	if _, err := os.Stat("DEV"); !os.IsNotExist(err) {
		log.SetLogLevel("debug")
		log.SetReportCaller(false)
	}

	godotenv.Load()

	consts.LoadEnvConfigs()

	info, err := utils.GetDBInfoFromEnv()
	if err != nil {
		log.Fatal("Error getting database info from env", "err", err)
	}

	err = db.ConnectDB(info)
	if err != nil {
		log.Fatal("Error connecting to database", "err", err)
	}
	defer db.CloseDB()

	ctx := context.Background()
	parseFlags(ctx)

	log.Info("Starting TRXd server")
	defer log.Info("Stopping TRXd server")

	go instancer.ReclaimLoop()

	app := api.SetupApp(ctx)
	err = app.Listen(":1337")
	if err != nil {
		log.Fatal("Error starting server", "err", err)
	}
}
