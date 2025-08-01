package tests

import (
	"context"
	"net/http"
	"testing"
	"trxd/api"
	"trxd/db"
	"trxd/utils"
	"trxd/utils/consts"
)

var testDeleteChallenge = []struct {
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
		testBody:         JSON{},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"chall_id": -1},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidChallengeID),
	},
	{
		testBody:       JSON{"chall_id": ""},
		expectedStatus: http.StatusOK,
	},
	{
		testBody:       JSON{"chall_id": ""},
		expectedStatus: http.StatusOK,
	},
}

func TestDeleteChallenge(t *testing.T) {
	db.DeleteAll()
	app := api.SetupApp()
	defer app.Shutdown()

	_, err := db.RegisterUser(context.Background(), "author", "author@test.test", "authorpass", db.UserRoleAuthor)
	if err != nil {
		t.Fatalf("Failed to register author user: %v", err)
	}

	var challID int32
	for _, test := range testDeleteChallenge {
		_, err := db.CreateCategory(context.Background(), "cat", "icon")
		if err != nil {
			t.Fatalf("Failed to create category: %v", err)
		}
		chall, err := db.CreateChallenge(context.Background(), "chall", "cat", "test-desc", db.DeployTypeNormal, 1, db.ScoreTypeStatic)
		if err != nil {
			t.Fatalf("Failed to create challenge: %v", err)
		}
		if chall != nil {
			challID = chall.ID
		}

		session := utils.NewApiTestSession(t, app)
		session.Post("/login", JSON{"email": "author@test.test", "password": "authorpass"}, http.StatusOK)
		if body, ok := test.testBody.(JSON); ok && body != nil {
			if content, ok := body["chall_id"]; ok && content == "" {
				test.testBody.(JSON)["chall_id"] = challID
			}
		}
		session.Delete("/challenge", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}
}
