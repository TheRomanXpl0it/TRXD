package main

import (
	"os"
	"trxd/db"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/joho/godotenv"
	"github.com/tde-nico/log"
)

// TODO: set store as redis + set configs (expire time, etc.)
var store = session.New()

func apiRegister(c *fiber.Ctx) error {
	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")

	if username == "" || email == "" || password == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Missing required fields")
	}
	if len(password) < 8 {
		return c.Status(fiber.StatusBadRequest).SendString("Password must be at least 8 characters long")
	}
	if len(username) > 64 || len(email) > 64 {
		return c.Status(fiber.StatusBadRequest).SendString("Username and email must not exceed 64 characters")
	}

	user, err := db.RegisterUser(username, email, password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error registering user")
	}
	if user == nil {
		return c.Status(fiber.StatusConflict).SendString("User already exists")
	}

	sess, err := store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error creating session")
	}

	sess.Set("username", user.Name)
	sess.Set("api-key", user.Apikey)
	sess.Save()

	return c.Status(fiber.StatusOK).SendString("User registered successfully")
}

func apiLogin(c *fiber.Ctx) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	if email == "" || password == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Missing required fields")
	}

	user, err := db.LoginUser(email, password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error logging in")
	}
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid email or password")
	}

	sess, err := store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error creating session")
	}

	sess.Set("username", user.Name)
	sess.Set("api-key", user.Apikey)
	sess.Save()

	return c.Status(fiber.StatusOK).SendString("User logged in successfully")
}

func apiLogout(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving session")
	}

	err = sess.Destroy()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error destroying session")
	}

	return c.Status(fiber.StatusOK).SendString("User logged out successfully")
}

func setupApp() *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "TRXd",
	})

	app.Post("/register", apiRegister)
	app.Post("/login", apiLogin)
	app.Post("/logout", apiLogout)

	return app
}

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

	app := setupApp()
	err = app.Listen(":1337")
	if err != nil {
		log.Fatal("Error starting server", "err", err)
	}
}
