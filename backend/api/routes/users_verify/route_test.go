package users_verify_test

import (
	"net/http"
	"testing"
	"trxd/api"
	"trxd/api/routes/users_register"
	"trxd/utils/consts"
	jwt_utils "trxd/utils/jwt"
	"trxd/utils/test_utils"

	"github.com/golang-jwt/jwt/v5"
)

type JSON map[string]interface{}

func errorf(val interface{}) JSON {
	return JSON{"error": val}
}

func TestMain(m *testing.M) {
	test_utils.Main(m)
}

func TestRoute(t *testing.T) {
	app := api.SetupApp(t.Context())
	defer api.Shutdown(app)

	session := test_utils.NewApiTestSession(t, app)
	session.Get("/verify", nil, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.InvalidToken))

	session.Get("/verify?token=AAA", nil, http.StatusUnauthorized)
	session.CheckResponse(errorf("token is malformed: token contains an invalid number of segments"))

	key := []byte("test-secret-key")
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt_utils.Map{"a": "b"})
	signed, err := token.SignedString(key)
	if err != nil {
		t.Fatalf("failed to sign JWT: %v", err)
	}
	session.Get("/verify?token="+signed, nil, http.StatusUnauthorized)
	session.CheckResponse(errorf("token is unverifiable: error while executing keyfunc: " + consts.InvalidSigningAlgorithm))

	tokenStr, err := jwt_utils.GenerateJWT(t.Context(), jwt_utils.Map{"a": "b"})
	if err != nil {
		t.Fatalf("failed to generate JWT: %v", err)
	}
	session.Get("/verify?token="+tokenStr, nil, http.StatusUnauthorized)
	session.CheckResponse(errorf(consts.InvalidToken))

	tokenStr, err = jwt_utils.GenerateJWT(t.Context(), jwt_utils.Map{"email": "m@m.m"})
	if err != nil {
		t.Fatalf("failed to generate JWT: %v", err)
	}
	session.Get("/verify?token="+tokenStr, nil, http.StatusUnauthorized)
	session.CheckResponse(errorf(consts.InvalidToken))

	ok, err := users_register.SetNXUserData(t.Context(), users_register.Data{
		Name:     "m",
		Email:    "m@m.m",
		Password: "12345678",
	})
	if err != nil {
		t.Fatalf("failed to set NX user data: %v", err)
	}
	if !ok {
		t.Fatal("expected to set NX user data successfully")
	}
	tokenStr, err = jwt_utils.GenerateJWT(t.Context(), jwt_utils.Map{"email": "m@m.m"})
	if err != nil {
		t.Fatalf("failed to generate JWT: %v", err)
	}
	session.Get("/verify?token="+tokenStr, nil, http.StatusOK)
	session.CheckResponse(nil)

	expected := JSON{
		"name":      "m",
		"role":      "Player",
		"team_id":   nil,
		"user_mode": false,
	}
	session.Get("/info", nil, http.StatusOK)
	body := session.Body()
	if body == nil {
		t.Fatal("Expected body")
	}
	test_utils.DeleteKeys(body, "id")
	test_utils.Compare(t, expected, body)

	ok, err = users_register.SetNXUserData(t.Context(), users_register.Data{
		Name:     "m",
		Email:    "m@m.m",
		Password: "12345678",
	})
	if err != nil {
		t.Fatalf("failed to set NX user data: %v", err)
	}
	if !ok {
		t.Fatal("expected to set NX user data successfully")
	}
	tokenStr, err = jwt_utils.GenerateJWT(t.Context(), jwt_utils.Map{"email": "m@m.m"})
	if err != nil {
		t.Fatalf("failed to generate JWT: %v", err)
	}
	session.Get("/verify?token="+tokenStr, nil, http.StatusConflict)
	session.CheckResponse(errorf(consts.UserAlreadyExists))

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "m@m.m", "password": "12345678"}, http.StatusOK)
	session.CheckResponse(nil)
}
