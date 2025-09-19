package tags_update_test

import (
	"fmt"
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
		testBody:         JSON{"old_name": "test"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"new_name": "test"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"chall_id": "", "old_name": strings.Repeat("a", consts.MaxTagNameLength+1), "new_name": "test-2"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf("OldName must not exceed 32"),
	},
	{
		testBody:         JSON{"chall_id": "", "old_name": "test", "new_name": strings.Repeat("a", consts.MaxTagNameLength+1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf("NewName must not exceed 32"),
	},
	{
		testBody:         JSON{"chall_id": -1, "old_name": "test", "new_name": "test-2"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf("ChallID must be at least 0"),
	},
	{
		testBody:       JSON{"chall_id": "", "old_name": "test", "new_name": "test-2"},
		expectedStatus: http.StatusOK,
	},
	{
		testBody:       JSON{"chall_id": "", "old_name": "test-2", "new_name": "test"},
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
	session.Post("/tags", JSON{"chall_id": chall.ID, "name": "test"}, http.StatusOK)

	for _, test := range testData {
		session := test_utils.NewApiTestSession(t, app)
		session.Post("/login", JSON{"email": "author@test.test", "password": "testpass"}, http.StatusOK)
		if body, ok := test.testBody.(JSON); ok && body != nil {
			if content, ok := body["chall_id"]; ok && content == "" {
				test.testBody.(JSON)["chall_id"] = chall.ID
			}
		}
		session.Patch("/tags", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "author@test.test", "password": "testpass"}, http.StatusOK)
	session.Get(fmt.Sprintf("/challenges/%d", chall.ID), nil, http.StatusOK)
	body := session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	tags := body.(map[string]interface{})["tags"].([]interface{})
	if len(tags) != 1 {
		t.Fatalf("Expected no tags, but got: %v", tags)
	}
	if tags[0].(string) != "test" {
		t.Fatalf("Expected tag to be 'test', but got: %v", tags[0])
	}
}
