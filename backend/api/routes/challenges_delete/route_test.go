package challenges_delete_test

import (
	"net/http"
	"testing"
	"trxd/api"
	"trxd/api/routes/categories_create"
	"trxd/api/routes/challenges_create"
	"trxd/db/sqlc"
	"trxd/utils/consts"
	"trxd/utils/test_utils"
)

type JSON map[string]interface{}

func errorf(val interface{}) JSON {
	return JSON{"error": val}
}

func TestMain(m *testing.M) {
	test_utils.Main(m, "../../../", "challenges_delete")
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

func TestRoute(t *testing.T) {
	app := api.SetupApp()
	defer app.Shutdown()

	test_utils.RegisterUser(t, "author", "author@test.test", "authorpass", sqlc.UserRoleAuthor)

	var challID int32
	for _, test := range testData {
		_, err := categories_create.CreateCategory(t.Context(), "cat", "icon")
		if err != nil {
			t.Fatalf("Failed to create category: %v", err)
		}
		chall, err := challenges_create.CreateChallenge(t.Context(), "chall", "cat", "test-desc", sqlc.DeployTypeNormal, 1, sqlc.ScoreTypeStatic)
		if err != nil {
			t.Fatalf("Failed to create challenge: %v", err)
		}
		if chall != nil {
			challID = chall.ID
		}

		session := test_utils.NewApiTestSession(t, app)
		session.Post("/users/login", JSON{"email": "author@test.test", "password": "authorpass"}, http.StatusOK)
		if body, ok := test.testBody.(JSON); ok && body != nil {
			if content, ok := body["chall_id"]; ok && content == "" {
				test.testBody.(JSON)["chall_id"] = challID
			}
		}
		session.Delete("/challenges/delete", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}
}
