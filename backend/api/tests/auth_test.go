package tests

import (
	"context"
	"net/http"
	"testing"
	"trxd/api"
	"trxd/db"
	"trxd/utils"
)

var testAuthMiddlewares = []struct {
	method           string
	endpoint         string
	expectedStatuses []int
}{
	{
		method:           "POST",
		endpoint:         "/login",
		expectedStatuses: []int{http.StatusBadRequest, http.StatusForbidden, http.StatusForbidden, http.StatusForbidden, http.StatusForbidden},
	},
	{
		method:           "GET",
		endpoint:         "/info",
		expectedStatuses: []int{http.StatusUnauthorized, http.StatusOK, http.StatusOK, http.StatusOK, http.StatusOK},
	},
	{
		method:           "POST",
		endpoint:         "/teams",
		expectedStatuses: []int{http.StatusUnauthorized, http.StatusForbidden, http.StatusBadRequest, http.StatusBadRequest, http.StatusBadRequest},
	},
	{
		method:           "POST",
		endpoint:         "/category",
		expectedStatuses: []int{http.StatusUnauthorized, http.StatusForbidden, http.StatusForbidden, http.StatusBadRequest, http.StatusBadRequest},
	},
	{
		method:           "PATCH",
		endpoint:         "/config",
		expectedStatuses: []int{http.StatusUnauthorized, http.StatusForbidden, http.StatusForbidden, http.StatusForbidden, http.StatusBadRequest},
	},
}

func TestAuthMiddlewares(t *testing.T) {
	db.DeleteAll()
	db.InitConfigs()
	app := api.SetupApp()
	defer app.Shutdown()

	err := db.UpdateConfig(context.Background(), "allow-register", "true")
	if err != nil {
		t.Fatalf("Failed to update config: %v", err)
	}

	users := [5]*db.User{}
	users[1], err = db.RegisterUser(context.Background(), "spectator", "spectator@test.test", "testpass", db.UserRoleSpectator)
	if err != nil {
		t.Fatalf("Failed to register spectator user: %v", err)
	}
	users[2], err = db.RegisterUser(context.Background(), "player", "player@test.test", "testpass", db.UserRolePlayer)
	if err != nil {
		t.Fatalf("Failed to register player user: %v", err)
	}
	users[3], err = db.RegisterUser(context.Background(), "author", "author@test.test", "testpass", db.UserRoleAuthor)
	if err != nil {
		t.Fatalf("Failed to register author user: %v", err)
	}
	users[4], err = db.RegisterUser(context.Background(), "admin", "admin@test.test", "testpass", db.UserRoleAdmin)
	if err != nil {
		t.Fatalf("Failed to register admin user: %v", err)
	}

	for _, test := range testAuthMiddlewares {
		for j, user := range users {
			session := utils.NewApiTestSession(t, app)
			if user != nil {
				session.Post("/login", JSON{"email": user.Email, "password": "testpass"}, http.StatusOK)
			}
			session.Request(test.method, test.endpoint, nil, test.expectedStatuses[j])
		}
	}
}
