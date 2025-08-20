package challenges_create_test

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
		testBody:       JSON{"name": "test", "category": "cat", "description": "test-desc", "type": "Normal", "max_points": 1, "score_type": "Static"},
		expectedStatus: http.StatusOK,
	},
	{
		testBody:         JSON{"name": "test", "category": "cat", "description": "test-desc", "type": "Normal", "max_points": 1, "score_type": "Static"},
		expectedStatus:   http.StatusConflict,
		expectedResponse: errorf(consts.ChallengeAlreadyExists),
	},
	{
		testBody:       JSON{"name": "test2", "category": "cat", "description": "test-desc", "type": "Normal", "max_points": 1, "score_type": "Static"},
		expectedStatus: http.StatusOK,
	},
}

func TestRoute(t *testing.T) {
	app := api.SetupApp()
	defer app.Shutdown()

	test_utils.RegisterUser(t, "author", "author@test.test", "authorpass", sqlc.UserRoleAuthor)
	session := test_utils.NewApiTestSession(t, app)
	session.Post("/users/login", JSON{"email": "author@test.test", "password": "authorpass"}, http.StatusOK)
	session.Post("/categories/create", JSON{"name": "cat", "icon": "icon"}, http.StatusOK)

	for _, test := range testData {
		session := test_utils.NewApiTestSession(t, app)
		session.Post("/users/login", JSON{"email": "author@test.test", "password": "authorpass"}, http.StatusOK)
		session.Post("/challenges/create", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}
}
