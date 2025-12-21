package challenges_create_test

import (
	"math"
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
		testBody:         JSON{"name": strings.Repeat("a", consts.MaxChallNameLen+1), "category": "cat", "description": "test-desc", "type": "Normal", "max_points": 1, "score_type": "Static"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(test_utils.Format(consts.MaxError, "Name", consts.MaxChallNameLen)),
	},
	{
		testBody:         JSON{"name": "test", "category": strings.Repeat("a", consts.MaxCategoryLen+1), "description": "test-desc", "type": "Normal", "max_points": 1, "score_type": "Static"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(test_utils.Format(consts.MaxError, "Category", consts.MaxCategoryLen)),
	},
	{
		testBody:         JSON{"name": "test", "category": "cat", "description": strings.Repeat("a", consts.MaxChallDescLen+1), "type": "Normal", "max_points": 1, "score_type": "Static"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(test_utils.Format(consts.MaxError, "Description", consts.MaxChallDescLen)),
	},
	{
		testBody:         JSON{"name": "test", "category": "cat", "description": "test-desc", "type": "aaaaa", "max_points": 1, "score_type": "Static"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(test_utils.Format(consts.OneOfError, "Type", strings.Join(consts.DeployTypesStr, " "))),
	},
	{
		testBody:         JSON{"name": "test", "category": "cat", "description": "test-desc", "type": "Normal", "max_points": -1, "score_type": "Static"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(test_utils.Format(consts.MinError, "MaxPoints", 0)),
	},
	{
		testBody:         JSON{"name": "test", "category": "cat", "description": "test-desc", "type": "Normal", "max_points": math.MaxInt32 + 1, "score_type": "Static"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidJSON),
	},
	{
		testBody:         JSON{"name": "test", "category": "cat", "description": "test-desc", "type": "Normal", "max_points": 1, "score_type": "aaaa"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(test_utils.Format(consts.OneOfError, "ScoreType", strings.Join(consts.ScoreTypesStr, " "))),
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
	app := api.SetupApp(t.Context())
	defer api.Shutdown(app)

	test_utils.RegisterUser(t, "author", "author@test.test", "authorpass", sqlc.UserRoleAuthor)
	session := test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "author@test.test", "password": "authorpass"}, http.StatusOK)
	session.Post("/categories", JSON{"name": "cat"}, http.StatusOK)

	for _, test := range testData {
		session := test_utils.NewApiTestSession(t, app)
		session.Post("/login", JSON{"email": "author@test.test", "password": "authorpass"}, http.StatusOK)
		session.Post("/challenges", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}
}
