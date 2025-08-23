package categories_update_test

import (
	"net/http"
	"strings"
	"testing"
	"trxd/api"
	"trxd/db/sqlc"
	"trxd/utils/consts"
	"trxd/utils/test_utils"
)

type JSON map[string]interface{}

func errorf(val interface{}) JSON {
	return JSON{"error": val}
}

func TestMain(m *testing.M) {
	test_utils.Main(m)
}

var testData = []struct {
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
		testBody:         JSON{"new_icon": "test"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"name": strings.Repeat("a", consts.MaxCategoryLength+1), "new_icon": "test"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.LongCategory),
	},
	{
		testBody:         JSON{"name": "test", "new_name": strings.Repeat("a", consts.MaxCategoryLength+1), "new_icon": "test"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.LongCategory),
	},
	{
		testBody:         JSON{"name": "test", "new_name": "test", "new_icon": strings.Repeat("a", consts.MaxIconLength+1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.LongIcon),
	},
	{
		testBody:         JSON{"name": "test", "new_icon": "AAA"},
		expectedStatus:   http.StatusNotFound,
		expectedResponse: errorf(consts.CategoryNotFound),
	},
	{
		testBody:       JSON{"name": "cat-1", "new_icon": "AAA"},
		expectedStatus: http.StatusOK,
	},
	{
		testBody:       JSON{"name": "cat-1", "new_name": "category-1"},
		expectedStatus: http.StatusOK,
	},
	{
		testBody:         JSON{"name": "cat-1", "new_name": "category-1"},
		expectedStatus:   http.StatusNotFound,
		expectedResponse: errorf(consts.CategoryNotFound),
	},
	{
		testBody:       JSON{"name": "category-1", "new_name": "challs-1", "new_icon": "BBB"},
		expectedStatus: http.StatusOK,
	},
}

func TestRoute(t *testing.T) {
	app := api.SetupApp()
	defer app.Shutdown()

	test_utils.RegisterUser(t, "author", "author@test.test", "authorpass", sqlc.UserRoleAuthor)

	for _, test := range testData {
		session := test_utils.NewApiTestSession(t, app)
		session.Post("/login", JSON{"email": "author@test.test", "password": "authorpass"}, http.StatusOK)
		session.Patch("/categories", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}

	session := test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "author@test.test", "password": "authorpass"}, http.StatusOK)
	session.Get("/challenges", nil, http.StatusOK)
	body := session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}

	count_cat_1 := 0
	count_challs_1 := 0
	for _, chall := range body.([]interface{}) {
		switch chall.(map[string]interface{})["category"] {
		case "cat-1":
			count_cat_1++
		case "challs-1":
			count_challs_1++
		}
	}

	if count_cat_1 != 0 && count_challs_1 != 3 {
		t.Fatalf("Unexpected challenge counts: cat-1: %d, challs-1: %d", count_cat_1, count_challs_1)
	}
}
