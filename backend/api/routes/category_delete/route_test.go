package category_delete_test

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"trxd/api"
	"trxd/api/routes/category_create"
	"trxd/api/routes/user_register"
	"trxd/db"
	"trxd/utils/consts"
	"trxd/utils/test_utils"
)

type JSON map[string]interface{}

func errorf(val interface{}) JSON {
	return JSON{"error": val}
}

func TestMain(m *testing.M) {
	test_utils.Main(m, "../../../", "category_delete")
}

var testCategoryDelete = []struct {
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

func TestDeleteCategoryDelete(t *testing.T) {
	app := api.SetupApp()
	defer app.Shutdown()

	_, err := user_register.RegisterUser(context.Background(), "author", "author@test.test", "authorpass", db.UserRoleAuthor)
	if err != nil {
		t.Fatalf("Failed to register author user: %v", err)
	}

	for _, test := range testCategoryDelete {
		_, err := category_create.CreateCategory(context.Background(), "cat", "icon")
		if err != nil {
			t.Fatalf("Failed to create category: %v", err)
		}

		session := test_utils.NewApiTestSession(t, app)
		session.Post("/login", JSON{"email": "author@test.test", "password": "authorpass"}, http.StatusOK)
		session.Delete("/category", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}
}
