package tags_create_test

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
		testBody:         JSON{"chall_id": ""},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"name": "test"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"chall_id": "", "name": strings.Repeat("a", consts.MaxTagNameLen+1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf("Name must not exceed 32"),
	},
	{
		testBody:         JSON{"chall_id": -1, "name": "test"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf("ChallID must be at least 0"),
	},
	{
		testBody:         JSON{"chall_id": 99999, "name": "test"},
		expectedStatus:   http.StatusNotFound,
		expectedResponse: errorf(consts.ChallengeNotFound),
	},
	{
		testBody:       JSON{"chall_id": "", "name": "test"},
		expectedStatus: http.StatusOK,
	},
	{
		testBody:         JSON{"chall_id": "", "name": "test"},
		expectedStatus:   http.StatusConflict,
		expectedResponse: errorf(consts.TagAlreadyExists),
	},
	{
		testBody:       JSON{"chall_id": "", "name": "test-2"},
		expectedStatus: http.StatusOK,
	},
}

func TestRoute(t *testing.T) {
	app := api.SetupApp()
	defer app.Shutdown()

	test_utils.RegisterUser(t, "author", "author@test.test", "testpass", sqlc.UserRoleAuthor)
	session := test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "author@test.test", "password": "testpass"}, http.StatusOK)
	session.Post("/categories", JSON{"name": "cat", "icon": "icon"}, http.StatusOK)
	chall := test_utils.CreateChallenge(t, "chall", "cat", "test-desc", sqlc.DeployTypeNormal, 1, sqlc.ScoreTypeStatic)

	for _, test := range testData {
		session := test_utils.NewApiTestSession(t, app)
		session.Post("/login", JSON{"email": "author@test.test", "password": "testpass"}, http.StatusOK)
		if body, ok := test.testBody.(JSON); ok && body != nil {
			if content, ok := body["chall_id"]; ok && content == "" {
				test.testBody.(JSON)["chall_id"] = chall.ID
			}
		}
		session.Post("/tags", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}
}
