package flags_submit_test

import (
	"net/http"
	"strings"
	"testing"
	"trxd/api"
	"trxd/api/routes/categories_create"
	"trxd/api/routes/challenges_create"
	"trxd/api/routes/flags_create"
	"trxd/api/routes/teams_register"
	"trxd/db/sqlc"
	"trxd/utils/consts"
	"trxd/utils/test_utils"
)

type JSON map[string]interface{}

func errorf(val interface{}) JSON {
	return JSON{"error": val}
}

func TestMain(m *testing.M) {
	test_utils.Main(m, "../../../", "flags_submit")
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
	app := api.SetupApp()
	defer app.Shutdown()

	session := test_utils.NewApiTestSession(t, app)
	session.Post("/users/register", JSON{"username": "test2", "email": "test2@test.test", "password": "testpass"}, http.StatusOK)
	session.Post("/flags/submit", JSON{"chall_id": 0, "flag": "flag{test}"}, http.StatusForbidden)
	session.CheckResponse(errorf(consts.Forbidden))

	test_utils.RegisterUser(t, "test3", "test3@test.test", "testpass", sqlc.UserRoleAdmin)
	session = test_utils.NewApiTestSession(t, app)
	session.Post("/users/login", JSON{"email": "test3@test.test", "password": "testpass"}, http.StatusOK)
	session.Post("/flags/submit", JSON{"chall_id": 0, "flag": "flag{test}"}, http.StatusNotFound)
	session.CheckResponse(errorf(consts.ChallengeNotFound))

	user := test_utils.RegisterUser(t, "test", "test@test.test", "testpass", sqlc.UserRolePlayer)
	team, err := teams_register.RegisterTeam(t.Context(), "test-team", "teampasswd", user.ID)
	if err != nil {
		t.Fatalf("Failed to register test team: %v", err)
	}
	if team == nil {
		t.Fatal("Team registration returned nil")
	}
	user2 := test_utils.RegisterUser(t, "test-2", "test-2@test.test", "testpass", sqlc.UserRolePlayer)
	team2, err := teams_register.RegisterTeam(t.Context(), "test-team-2", "teampasswd", user2.ID)
	if err != nil {
		t.Fatalf("Failed to register test team 2: %v", err)
	}
	if team2 == nil {
		t.Fatal("Team2 registration returned nil")
	}

	cat, err := categories_create.CreateCategory(t.Context(), "cat", "icon")
	if err != nil {
		t.Fatalf("Failed to create category: %v", err)
	}
	if cat == nil {
		t.Fatal("Category creation returned nil")
	}
	chall, err := challenges_create.CreateChallenge(t.Context(), "chall", cat.Name, "test-desc", sqlc.DeployTypeNormal, 1, sqlc.ScoreTypeDynamic)
	if err != nil {
		t.Fatalf("Failed to create challenge: %v", err)
	}
	if chall == nil {
		t.Fatal("Challenge creation returned nil")
	}
	flag, err := flags_create.CreateFlag(t.Context(), chall.ID, "flag{test}", false)
	if err != nil {
		t.Fatalf("Failed to create flag: %v", err)
	}
	if flag == nil {
		t.Fatal("Flag creation returned nil")
	}

	for _, test := range testData {
		session := test_utils.NewApiTestSession(t, app)
		if test.secondUser {
			session.Post("/users/login", JSON{"email": "test-2@test.test", "password": "testpass"}, http.StatusOK)
		} else {
			session.Post("/users/login", JSON{"email": "test@test.test", "password": "testpass"}, http.StatusOK)
		}
		if body, ok := test.testBody.(JSON); ok && body != nil {
			if content, ok := body["chall_id"]; ok && content == "" {
				test.testBody.(JSON)["chall_id"] = chall.ID
			}
		}
		session.Post("/flags/submit", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}
}
