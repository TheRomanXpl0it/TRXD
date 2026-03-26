package teams_join_get_test

import (
	"fmt"
	"net/http"
	"testing"
	"trxd/api"
	"trxd/utils/consts"
	"trxd/utils/jwt"
	"trxd/utils/test_utils"
)

type JSON map[string]any

func errorf(val any) JSON {
	return JSON{"error": val}
}

func Json(val any) map[string]any {
	return val.(map[string]any)
}

func Int32(val any) int32 {
	return int32(val.(float64))
}

func TestMain(m *testing.M) {
	test_utils.Main(m)
}

func GenerateToken(t *testing.T, data jwt.Map) string {
	token, err := jwt.GenerateJWT(t.Context(), data)
	if err != nil {
		t.Fatalf("Failed to generate JWT (%+v): %v", data, err)
	}

	return token
}

func TestRoute(t *testing.T) {
	app := api.SetupApp(t.Context())
	defer api.Shutdown(app)

	// Get Token

	session := test_utils.NewApiTestSession(t, app)
	session.Post("/register", JSON{"name": "test", "email": "test@test.test", "password": "testpass"}, http.StatusOK)
	session.CheckResponse(nil)

	session.Get("/teams/join", nil, http.StatusNotFound)
	session.CheckResponse(errorf(consts.TeamNotFound))

	session.Post("/teams/register", JSON{"name": "test", "password": "testpass"}, http.StatusOK)
	session.CheckResponse(nil)

	session.Get("/info", nil, http.StatusOK)
	body := session.Body()
	tid := Int32(Json(body)["team_id"])

	session.Get("/teams/join", nil, http.StatusOK)
	body = session.Body()
	token := Json(body)["token"].(string)
	jwtMap, err := jwt.ParseAndValidateJWT(t.Context(), token)
	if err != nil {
		t.Fatalf("Failed to parse and validate JWT: %v", err)
	}

	if Int32(jwtMap["team_id"]) != tid {
		t.Fatalf("Expected tid %v in JWT, got %v", tid, Int32(jwtMap["team_id"]))
	}

	// Verify & Join

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/register", JSON{"name": "test2", "email": "test2@test.test", "password": "testpass"}, http.StatusOK)
	session.CheckResponse(nil)

	invalidToken := "invalidtoken"
	session.Get(fmt.Sprintf("/teams/join?token=%s", invalidToken), nil, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.InvalidToken))

	missingTID := GenerateToken(t, jwt.Map{"a": "b"})
	session.Get(fmt.Sprintf("/teams/join?token=%s", missingTID), nil, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.InvalidToken))

	tidNotANumber := GenerateToken(t, jwt.Map{"team_id": "notanumber"})
	session.Get(fmt.Sprintf("/teams/join?token=%s", tidNotANumber), nil, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.InvalidToken))

	team404 := GenerateToken(t, jwt.Map{"team_id": 9999999})
	session.Get(fmt.Sprintf("/teams/join?token=%s", team404), nil, http.StatusNotFound)
	session.CheckResponse(errorf(consts.TeamNotFound))

	session.Get(fmt.Sprintf("/teams/join?token=%s", token), nil, http.StatusOK)
	session.CheckResponse(nil)

	session.Get("/info", nil, http.StatusOK)
	body = session.Body()
	tid2 := Int32(Json(body)["team_id"])

	if tid != tid2 {
		t.Fatalf("Expected tid %v after joining team, got %v", tid, tid2)
	}
}
