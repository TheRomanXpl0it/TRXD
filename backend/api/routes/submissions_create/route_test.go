package submissions_create_test

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
	secondUser       bool
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
		testBody:         JSON{"chall_id": "", "flag": strings.Repeat("a", consts.MaxFlagLen+1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf("Flag must not exceed 128"),
	},
	{
		testBody:         JSON{"chall_id": -1, "flag": "flag{test}"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf("ChallID must be at least 0"),
	},
	{
		testBody:         JSON{"chall_id": 99999, "flag": "flag{test}"},
		expectedStatus:   http.StatusNotFound,
		expectedResponse: errorf(consts.ChallengeNotFound),
	},
	{
		testBody:         JSON{"chall_id": math.MaxInt32 + 1, "flag": "flag{test}"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidJSON),
	},
	{
		testBody:         JSON{"chall_id": "", "flag": "test"},
		expectedStatus:   http.StatusOK,
		expectedResponse: JSON{"status": sqlc.SubmissionStatusWrong, "first_blood": false},
	},
	{
		testBody:         JSON{"chall_id": "", "flag": "flag{test}"},
		expectedStatus:   http.StatusOK,
		expectedResponse: JSON{"status": sqlc.SubmissionStatusCorrect, "first_blood": true},
	},
	{
		testBody:         JSON{"chall_id": "", "flag": "flag{test}"},
		expectedStatus:   http.StatusOK,
		expectedResponse: JSON{"status": sqlc.SubmissionStatusRepeated, "first_blood": false},
	},
	{
		testBody:         JSON{"chall_id": "", "flag": " flag{test} "},
		expectedStatus:   http.StatusOK,
		expectedResponse: JSON{"status": sqlc.SubmissionStatusRepeated, "first_blood": false},
	},
	{
		testBody:         JSON{"chall_id": "", "flag": "test"},
		expectedStatus:   http.StatusOK,
		expectedResponse: JSON{"status": sqlc.SubmissionStatusWrong, "first_blood": false},
		secondUser:       true,
	},
	{
		testBody:         JSON{"chall_id": "", "flag": "flag{test}"},
		expectedStatus:   http.StatusOK,
		expectedResponse: JSON{"status": sqlc.SubmissionStatusCorrect, "first_blood": false},
		secondUser:       true,
	},
	{
		testBody:         JSON{"chall_id": "", "flag": "flag{test}"},
		expectedStatus:   http.StatusOK,
		expectedResponse: JSON{"status": sqlc.SubmissionStatusRepeated, "first_blood": false},
		secondUser:       true,
	},
}

func TestRoute(t *testing.T) {
	app := api.SetupApp(t.Context())
	defer app.Shutdown()

	session := test_utils.NewApiTestSession(t, app)
	session.Post("/register", JSON{"name": "test2", "email": "test2@test.test", "password": "testpass"}, http.StatusOK)
	session.Post("/submissions", JSON{"chall_id": 0, "flag": "flag{test}"}, http.StatusForbidden)
	session.CheckResponse(errorf(consts.Forbidden))

	test_utils.RegisterUser(t, "test3", "test3@test.test", "testpass", sqlc.UserRoleAdmin)
	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "test3@test.test", "password": "testpass"}, http.StatusOK)
	session.Post("/submissions", JSON{"chall_id": 0, "flag": "flag{test}"}, http.StatusNotFound)
	session.CheckResponse(errorf(consts.ChallengeNotFound))

	test_utils.RegisterUser(t, "test", "test@test.test", "testpass", sqlc.UserRolePlayer)
	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "test@test.test", "password": "testpass"}, http.StatusOK)
	session.Post("/teams/register", JSON{"name": "test-team", "password": "teampasswd"}, http.StatusOK)

	test_utils.RegisterUser(t, "test-2", "test-2@test.test", "testpass", sqlc.UserRolePlayer)
	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "test-2@test.test", "password": "testpass"}, http.StatusOK)
	session.Post("/teams/register", JSON{"name": "test-team-2", "password": "teampasswd"}, http.StatusOK)

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "test3@test.test", "password": "testpass"}, http.StatusOK)
	session.Post("/categories", JSON{"name": "cat", "icon": "icon"}, http.StatusOK)
	chall := test_utils.CreateChallenge(t, "chall", "cat", "test-desc", sqlc.DeployTypeNormal, 1, sqlc.ScoreTypeDynamic)
	test_utils.UnveilChallenge(t, chall.ID)
	session.Post("/flags", JSON{"chall_id": chall.ID, "flag": "flag{test}", "regex": false}, http.StatusOK)

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "test-2@test.test", "password": "testpass"}, http.StatusOK)
	session.Get("/challenges", nil, http.StatusOK)
	body := session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	var challID5 int32
	for _, chall := range body.([]interface{}) {
		if chall.(map[string]interface{})["name"] == "chall-5" {
			challID5 = int32(chall.(map[string]interface{})["id"].(float64))
		}
	}
	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "test-2@test.test", "password": "testpass"}, http.StatusOK)
	session.Post("/submissions", JSON{"chall_id": challID5, "flag": "flag{test}"}, http.StatusNotFound)
	session.CheckResponse(errorf(consts.ChallengeNotFound))

	for _, test := range testData {
		session := test_utils.NewApiTestSession(t, app)
		if test.secondUser {
			session.Post("/login", JSON{"email": "test-2@test.test", "password": "testpass"}, http.StatusOK)
		} else {
			session.Post("/login", JSON{"email": "test@test.test", "password": "testpass"}, http.StatusOK)
		}
		if body, ok := test.testBody.(JSON); ok && body != nil {
			if content, ok := body["chall_id"]; ok && content == "" {
				test.testBody.(JSON)["chall_id"] = chall.ID
			}
		}
		session.Post("/submissions", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}
}
