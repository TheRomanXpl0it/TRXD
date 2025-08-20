package flags_create_test

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
		testBody:         JSON{"flag": "test"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"chall_id": "", "flag": strings.Repeat("a", consts.MaxFlagLength+1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.LongFlag),
	},
	{
		testBody:         JSON{"chall_id": 99999, "flag": "flag{test}"},
		expectedStatus:   http.StatusNotFound,
		expectedResponse: errorf(consts.ChallengeNotFound),
	},
	{
		testBody:       JSON{"chall_id": "", "flag": "test"},
		expectedStatus: http.StatusOK,
	},
	{
		testBody:       JSON{"chall_id": "", "flag": `flag\{test\}`, "regex": true},
		expectedStatus: http.StatusOK,
	},
	{
		testBody:         JSON{"chall_id": "", "flag": "test"},
		expectedStatus:   http.StatusConflict,
		expectedResponse: errorf(consts.FlagAlreadyExists),
	},
}

func TestRoute(t *testing.T) {
	app := api.SetupApp()
	defer app.Shutdown()

	test_utils.RegisterUser(t, "test", "test@test.test", "testpass", sqlc.UserRoleAuthor)
	session := test_utils.NewApiTestSession(t, app)
	session.Post("/users/login", JSON{"email": "test@test.test", "password": "testpass"}, http.StatusOK)
	session.Post("/categories/create", JSON{"name": "cat", "icon": "icon"}, http.StatusOK)
	chall := test_utils.CreateChallenge(t, "chall", "cat", "test-desc", sqlc.DeployTypeNormal, 1, sqlc.ScoreTypeStatic)

	for _, test := range testData {
		session := test_utils.NewApiTestSession(t, app)
		session.Post("/users/login", JSON{"email": "test@test.test", "password": "testpass"}, http.StatusOK)
		if body, ok := test.testBody.(JSON); ok && body != nil {
			if content, ok := body["chall_id"]; ok && content == "" {
				test.testBody.(JSON)["chall_id"] = chall.ID
			}
		}
		session.Post("/flags/create", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}
}
