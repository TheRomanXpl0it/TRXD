package challenge_delete_test

import (
	"net/http"
	"testing"
	"trxd/api"
	"trxd/api/routes/category_create"
	"trxd/api/routes/challenge_create"
	"trxd/api/routes/user_register"
	"trxd/db/sqlc"
	"trxd/utils/consts"
	"trxd/utils/test_utils"
)

type JSON map[string]interface{}

func errorf(val interface{}) JSON {
	return JSON{"error": val}
}

func TestMain(m *testing.M) {
	test_utils.Main(m, "../../../", "challenge_delete")
}

var testChallengeDelete = []struct {
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

func TestChallengeDelete(t *testing.T) {
	app := api.SetupApp()
	defer app.Shutdown()

	_, err := user_register.RegisterUser(t.Context(), "author", "author@test.test", "authorpass", sqlc.UserRoleAuthor)
	if err != nil {
		t.Fatalf("Failed to register author user: %v", err)
	}

	var challID int32
	for _, test := range testChallengeDelete {
		_, err := category_create.CreateCategory(t.Context(), "cat", "icon")
		if err != nil {
			t.Fatalf("Failed to create category: %v", err)
		}
		chall, err := challenge_create.CreateChallenge(t.Context(), "chall", "cat", "test-desc", sqlc.DeployTypeNormal, 1, sqlc.ScoreTypeStatic)
		if err != nil {
			t.Fatalf("Failed to create challenge: %v", err)
		}
		if chall != nil {
			challID = chall.ID
		}

		session := test_utils.NewApiTestSession(t, app)
		session.Post("/login", JSON{"email": "author@test.test", "password": "authorpass"}, http.StatusOK)
		if body, ok := test.testBody.(JSON); ok && body != nil {
			if content, ok := body["chall_id"]; ok && content == "" {
				test.testBody.(JSON)["chall_id"] = challID
			}
		}
		session.Delete("/challenges", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}
}
