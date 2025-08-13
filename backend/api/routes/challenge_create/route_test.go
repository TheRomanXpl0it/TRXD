package challenge_create_test

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
	test_utils.Main(m, "../../../", "challenge_create")
}

var testChallengeCreate = []struct {
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
		testBody:         JSON{"category": "test"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"description": "test"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"type": "test"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"max_points": 1},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"score_type": "test"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"name": strings.Repeat("a", consts.MaxChallNameLength+1), "category": "cat", "description": "test-desc", "type": "Normal", "max_points": 1, "score_type": "Static"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.LongChallName),
	},
	{
		testBody:         JSON{"name": "test", "category": strings.Repeat("a", consts.MaxCategoryLength+1), "description": "test-desc", "type": "Normal", "max_points": 1, "score_type": "Static"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.LongCategory),
	},
	{
		testBody:         JSON{"name": "test", "category": "cat", "description": strings.Repeat("a", consts.MaxChallDescLength+1), "type": "Normal", "max_points": 1, "score_type": "Static"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.LongChallDesc),
	},
	{
		testBody:         JSON{"name": "test", "category": "cat", "description": "test-desc", "type": "aaaaa", "max_points": 1, "score_type": "Static"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidChallType),
	},
	{
		testBody:         JSON{"name": "test", "category": "cat", "description": "test-desc", "type": "Normal", "max_points": 0, "score_type": "Static"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidChallMaxPoints),
	},
	{
		testBody:         JSON{"name": "test", "category": "cat", "description": "test-desc", "type": "Normal", "max_points": 1, "score_type": "aaaa"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidChallScoreType),
	},
	{
		testBody:         JSON{"name": "test3", "category": "cat2", "description": "test-desc", "type": "Normal", "max_points": 1, "score_type": "Static"},
		expectedStatus:   http.StatusNotFound,
		expectedResponse: errorf(consts.CategoryNotFound),
	},
	{
		testBody:         JSON{"name": "test", "category": "cat", "description": "test-desc", "type": "Normal", "max_points": 1, "score_type": "Static"},
		expectedStatus:   http.StatusOK,
		expectedResponse: JSON{"name": "test", "category": "cat"},
	},
	{
		testBody:         JSON{"name": "test", "category": "cat", "description": "test-desc", "type": "Normal", "max_points": 1, "score_type": "Static"},
		expectedStatus:   http.StatusConflict,
		expectedResponse: errorf(consts.ChallengeAlreadyExists),
	},
	{
		testBody:         JSON{"name": "test2", "category": "cat", "description": "test-desc", "type": "Normal", "max_points": 1, "score_type": "Static"},
		expectedStatus:   http.StatusOK,
		expectedResponse: JSON{"name": "test2", "category": "cat"},
	},
}

func TestChallengeCreate(t *testing.T) {
	app := api.SetupApp()
	defer app.Shutdown()

	_, err := user_register.RegisterUser(context.Background(), "author", "author@test.test", "authorpass", db.UserRoleAuthor)
	if err != nil {
		t.Fatalf("Failed to register author user: %v", err)
	}
	_, err = category_create.CreateCategory(context.Background(), "cat", "icon")
	if err != nil {
		t.Fatalf("Failed to create category: %v", err)
	}

	for _, test := range testChallengeCreate {
		session := test_utils.NewApiTestSession(t, app)
		session.Post("/login", JSON{"email": "author@test.test", "password": "authorpass"}, http.StatusOK)
		session.Post("/challenges", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}
}
