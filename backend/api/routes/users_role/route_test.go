package users_role_test

import (
	"fmt"
	"math"
	"net/http"
	"testing"
	"trxd/api"
	"trxd/utils/consts"
	"trxd/utils/test_utils"
)

type JSON map[string]any

func errorf(val any) JSON {
	return JSON{"error": val}
}

func Json(val any) map[string]any {
	return val.(map[string]any)
}

func List(val any) []any {
	return val.([]any)
}

func Int32(val any) int32 {
	return int32(val.(float64))
}

func Int(val any) int {
	return int(val.(float64))
}

func TestMain(m *testing.M) {
	test_utils.Main(m)
}

var testData = []struct {
	testBody         any
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
		testBody:         JSON{"user_id": 0},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"new_role": "Admin"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"user_id": -1, "new_role": "Admin"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(test_utils.Format(consts.MinError, "UserID", 0)),
	},
	{
		testBody:         JSON{"user_id": math.MaxInt32 + 1, "new_role": "Admin"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidJSON),
	},
	{
		testBody:         JSON{"user_id": 0, "new_role": "aaa"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(test_utils.Format(consts.OneOfError, "NewRole", consts.RolesStr)),
	},
	{
		testBody:         JSON{"user_id": 0, "new_role": "Admin"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidRole),
	},
}

func TestRoute(t *testing.T) {
	app := api.SetupApp(t.Context())
	defer api.Shutdown(app)

	session := test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "admin@email.com", "password": "testpass"}, http.StatusOK)
	session.CheckResponse(nil)
	session.Get("/users", nil, http.StatusOK)
	body := session.Body()
	var uid int32
	for _, user := range List(Json(body)["users"]) {
		if Json(user)["name"] == "a" {
			uid = Int32(Json(user)["id"])
		}
	}

	for _, test := range testData {
		if body, ok := test.testBody.(JSON); ok && body != nil {
			if content, ok := body["user_id"]; ok && content == 0 {
				test.testBody.(JSON)["user_id"] = uid
			}
		}
		session.Patch("/users/role", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}

	session.Patch("/users/role", JSON{"user_id": uid, "new_role": "Author"}, http.StatusOK)
	session.CheckResponse(nil)
	session.Get(fmt.Sprintf("/users/%d", uid), nil, http.StatusOK)
	body = session.Body()
	if Json(body)["role"] != "Author" {
		t.Fatalf("Expected role to be 'Author', got '%v'", Json(body)["role"])
	}
	if Int(Json(body)["score"]) != 0 {
		t.Fatalf("Expected score to be 0, got %v", Json(body)["score"])
	}

	session.Patch("/users/role", JSON{"user_id": uid, "new_role": "Spectator"}, http.StatusOK)
	session.CheckResponse(nil)
	session.Get(fmt.Sprintf("/users/%d", uid), nil, http.StatusOK)
	body = session.Body()
	if Json(body)["role"] != "Spectator" {
		t.Fatalf("Expected role to be 'Spectator', got '%v'", Json(body)["role"])
	}
	if Int(Json(body)["score"]) != 0 {
		t.Fatalf("Expected score to be 0, got %v", Json(body)["score"])
	}

	session.Patch("/users/role", JSON{"user_id": uid, "new_role": "Player"}, http.StatusOK)
	session.CheckResponse(nil)
	session.Get(fmt.Sprintf("/users/%d", uid), nil, http.StatusOK)
	body = session.Body()
	if Json(body)["role"] != "Player" {
		t.Fatalf("Expected role to be 'Player', got '%v'", Json(body)["role"])
	}
	if Int(Json(body)["score"]) == 0 {
		t.Fatalf("Expected score to be restored, got %v", Json(body)["score"])
	}
}
