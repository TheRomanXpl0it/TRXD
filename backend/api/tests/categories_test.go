package tests

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"trxd/api"
	"trxd/db"
	"trxd/utils"
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

	_, err := db.RegisterUser(context.Background(), "author", "author@test.test", "authorpass", db.UserRoleAuthor)
	if err != nil {
		t.Fatalf("Failed to register author user: %v", err)
	}

	for _, test := range testCreateCategory {
		session := utils.NewApiTestSession(t, app)
		session.Post("/api/login", JSON{"email": "author@test.test", "password": "authorpass"}, http.StatusOK)
		session.Post("/api/author/category", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}
}

var testDeleteCategory = []struct {
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
		testBody:         JSON{},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"category": strings.Repeat("a", consts.MaxCategoryLength+1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.LongCategory),
	},
	{
		testBody:       JSON{"category": "cat"},
		expectedStatus: http.StatusOK,
	},
	{
		testBody:       JSON{"category": "cat"},
		expectedStatus: http.StatusOK,
	},
}

func TestDeleteCategory(t *testing.T) {
	db.DeleteAll()
	app := api.SetupApp()
	defer app.Shutdown()

	_, err := db.RegisterUser(context.Background(), "author", "author@test.test", "authorpass", db.UserRoleAuthor)
	if err != nil {
		t.Fatalf("Failed to register author user: %v", err)
	}

	for _, test := range testDeleteCategory {
		_, err := db.CreateCategory(context.Background(), "cat", "icon")
		if err != nil {
			t.Fatalf("Failed to create category: %v", err)
		}

		session := utils.NewApiTestSession(t, app)
		session.Post("/api/login", JSON{"email": "author@test.test", "password": "authorpass"}, http.StatusOK)
		session.Delete("/api/author/category", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}
}
