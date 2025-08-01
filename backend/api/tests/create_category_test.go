package tests

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"trxd/api"
	"trxd/db"
	"trxd/utils/consts"
)

var testCreateCategory = []struct {
	testBody         interface{}
	expectedStatus   int
	expectedResponse JSON
}{
	{
		testBody:         nil,
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidJSON),
	},
	{
		testBody:         JSON{"name": "test"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"icon": "test"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"name": strings.Repeat("a", consts.MaxCategoryLength+1), "icon": "test"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.LongCategory),
	},
	{
		testBody:         JSON{"name": "test", "icon": strings.Repeat("a", consts.MaxIconLength+1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.LongIcon),
	},
	{
		testBody:         JSON{"name": "test", "icon": "test"},
		expectedStatus:   http.StatusOK,
		expectedResponse: JSON{"name": "test", "icon": "test"},
	},
	{
		testBody:         JSON{"name": "test", "icon": "test"},
		expectedStatus:   http.StatusConflict,
		expectedResponse: errorf(consts.CategoryAlreadyExists),
	},
	{
		testBody:         JSON{"name": "test2", "icon": "test"},
		expectedStatus:   http.StatusOK,
		expectedResponse: JSON{"name": "test2", "icon": "test"},
	},
}

func TestCreateCategory(t *testing.T) {
	db.DeleteAll()
	app := api.SetupApp()
	defer app.Shutdown()

	_, err := db.RegisterUser(context.Background(), "admin", "admin@test.test", "adminpass", db.UserRoleAdmin)
	if err != nil {
		t.Fatalf("Failed to register admin user: %v", err)
	}

	for _, test := range testCreateCategory {
		session := newApiTestSession(t, app)
		session.Post("/login", JSON{"email": "admin@test.test", "password": "adminpass"}, http.StatusOK)
		session.Post("/create-category", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}
}
