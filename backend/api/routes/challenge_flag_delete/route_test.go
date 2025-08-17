package challenge_flag_delete_test

import (
	"net/http"
	"strings"
	"testing"
	"trxd/api"
	"trxd/api/routes/category_create"
	"trxd/api/routes/challenge_create"
	"trxd/api/routes/challenge_flag_create"
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
	test_utils.Main(m, "../../../", "flag_delete")
}

var testFlagDelete = []struct {
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
		testBody:       JSON{"chall_id": "", "flag": "test"},
		expectedStatus: http.StatusOK,
	},
}

func TestFlagDelete(t *testing.T) {
	app := api.SetupApp()
	defer app.Shutdown()

	user, err := user_register.RegisterUser(t.Context(), "test", "test@test.test", "testpass", sqlc.UserRoleAuthor)
	if err != nil {
		t.Fatalf("Failed to register author user: %v", err)
	}
	if user == nil {
		t.Fatal("User registration returned nil")
	}

	cat, err := category_create.CreateCategory(t.Context(), "cat", "icon")
	if err != nil {
		t.Fatalf("Failed to create category: %v", err)
	}
	if cat == nil {
		t.Fatal("Category creation returned nil")
	}
	chall, err := challenge_create.CreateChallenge(t.Context(), "chall", cat.Name, "test-desc", sqlc.DeployTypeNormal, 1, sqlc.ScoreTypeStatic)
	if err != nil {
		t.Fatalf("Failed to create challenge: %v", err)
	}
	if chall == nil {
		t.Fatal("Challenge creation returned nil")
	}

	for _, test := range testFlagDelete {
		_, err := challenge_flag_create.CreateFlag(t.Context(), chall.ID, "test", false)
		if err != nil {
			t.Fatalf("Failed to create flag: %v", err)
		}

		session := test_utils.NewApiTestSession(t, app)
		session.Post("/login", JSON{"email": "test@test.test", "password": "testpass"}, http.StatusOK)
		if body, ok := test.testBody.(JSON); ok && body != nil {
			if content, ok := body["chall_id"]; ok && content == "" {
				test.testBody.(JSON)["chall_id"] = chall.ID
			}
		}
		session.Delete("/flag", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}
}
