package tests

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"trxd/api"
	"trxd/db"
	"trxd/utils"
	"trxd/utils/consts"
)

var testCreateFlag = []struct {
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

func TestCreateFlag(t *testing.T) {
	db.DeleteAll()
	app := api.SetupApp()
	defer app.Shutdown()

	user, err := db.RegisterUser(context.Background(), "test", "test@test.test", "testpass", db.UserRoleAuthor)
	if err != nil {
		t.Fatalf("Failed to register test user: %v", err)
	}
	if user == nil {
		t.Fatal("User registration returned nil")
	}

	cat, err := db.CreateCategory(context.Background(), "cat", "icon")
	if err != nil {
		t.Fatalf("Failed to create category: %v", err)
	}
	if cat == nil {
		t.Fatal("Category creation returned nil")
	}
	chall, err := db.CreateChallenge(context.Background(), "chall", cat.Name.(string), "test-desc", db.DeployTypeNormal, 1, db.ScoreTypeStatic)
	if err != nil {
		t.Fatalf("Failed to create challenge: %v", err)
	}
	if chall == nil {
		t.Fatal("Challenge creation returned nil")
	}

	for _, test := range testCreateFlag {
		session := utils.NewApiTestSession(t, app)
		session.Post("/login", JSON{"email": "test@test.test", "password": "testpass"}, http.StatusOK)
		if body, ok := test.testBody.(JSON); ok && body != nil {
			if content, ok := body["chall_id"]; ok && content == "" {
				test.testBody.(JSON)["chall_id"] = chall.ID
			}
		}
		session.Post("/create-flag", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}
}
