package category_create_test

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"trxd/api"
	"trxd/api/routes/user_register"
	"trxd/db/sqlc"
	"trxd/utils/consts"
	"trxd/utils/test_utils"
)

type JSON map[string]interface{}

func errorf(val interface{}) JSON {
	return JSON{"error": val}
}

func TestMain(m *testing.M) {
	test_utils.Main(m, "../../../", "category_create")
}

var testCategoryCreate = []struct {
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

func TestCategoryCreate(t *testing.T) {
	app := api.SetupApp()
	defer app.Shutdown()

	_, err := user_register.RegisterUser(context.Background(), "author", "author@test.test", "authorpass", sqlc.UserRoleAuthor)
	if err != nil {
		t.Fatalf("Failed to register author user: %v", err)
	}

	for _, test := range testCategoryCreate {
		session := test_utils.NewApiTestSession(t, app)
		session.Post("/login", JSON{"email": "author@test.test", "password": "authorpass"}, http.StatusOK)
		session.Post("/category", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}
}
