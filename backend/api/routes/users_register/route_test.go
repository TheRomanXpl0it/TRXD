package users_register_test

import (
	"net/http"
	"strings"
	"testing"
	"trxd/api"
	"trxd/utils/consts"
	jwt_utils "trxd/utils/jwt"
	"trxd/utils/test_utils"

	"github.com/golang-jwt/jwt/v5"
)

type JSON map[string]any

func errorf(val any) JSON {
	return JSON{"error": val}
}

func Json(val any) map[string]any {
	return val.(map[string]any)
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
		testBody:         JSON{"name": "test"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"email": "test@test.test"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"password": "testpass"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"name": "test", "email": "test@test.test"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"name": "test", "password": "testpass"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"email": "test@test.test", "password": "testpass"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"name": "test", "email": "test@test.test", "password": strings.Repeat("a", consts.MinPasswordLen-1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(test_utils.Format(consts.MinError, "Password", consts.MinPasswordLen)),
	},
	{
		testBody:         JSON{"name": "test", "email": "test@test.test", "password": strings.Repeat("a", consts.MaxPasswordLen+1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(test_utils.Format(consts.MaxError, "Password", consts.MaxPasswordLen)),
	},
	{
		testBody:         JSON{"name": strings.Repeat("a", consts.MaxUserNameLen+1), "email": "test@test.test", "password": "testpass"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(test_utils.Format(consts.MaxError, "Name", consts.MaxUserNameLen)),
	},
	{
		testBody:         JSON{"name": "test", "email": strings.Repeat("a", consts.MaxEmailLen+1), "password": "testpass"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(test_utils.Format(consts.MaxError, "Email", consts.MaxEmailLen)),
	},
	{
		testBody:         JSON{"name": "test", "email": "invalid-email", "password": "testpass"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidEmail),
	},
	{
		testBody:       JSON{"name": "test", "email": "test@test.test", "password": "testpass"},
		expectedStatus: http.StatusOK,
	},
	{
		testBody:         JSON{"name": "test", "email": "test@test.test", "password": "testpass"},
		expectedStatus:   http.StatusConflict,
		expectedResponse: errorf(consts.UserAlreadyExists),
	},
	{
		testBody:         JSON{"name": "test", "email": "test1@test.test", "password": "testpass"},
		expectedStatus:   http.StatusConflict,
		expectedResponse: errorf(consts.UserAlreadyExists),
	},
	{
		testBody:       JSON{"name": "test2", "email": "test1@test.test", "password": "testpass"},
		expectedStatus: http.StatusOK,
	},
}

func TestRoute(t *testing.T) {
	app := api.SetupApp(t.Context())
	defer api.Shutdown(app)

	test_utils.UpdateConfig(t, "allow-register", "false")
	session := test_utils.NewApiTestSession(t, app)
	session.Post("/register", JSON{"name": "test-user", "email": "allow@test.test", "password": "testpass"}, http.StatusForbidden)
	session.CheckResponse(errorf(consts.DisabledRegistrations))

	test_utils.UpdateConfig(t, "allow-register", "true")
	session = test_utils.NewApiTestSession(t, app)
	session.Post("/register", JSON{"name": "test-user", "email": "allow@test.test", "password": "testpass"}, http.StatusOK)
	session.CheckResponse(nil)
	session.Post("/register", JSON{"name": "test-user-1", "email": "allow+1@test.test", "password": "testpass"}, http.StatusForbidden)
	session.CheckResponse(errorf(consts.AlreadyRegistered))

	for _, test := range testData {
		session := test_utils.NewApiTestSession(t, app)
		session.Post("/register", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}

	// APP 2

	app2 := api.SetupApp(t.Context())
	defer api.Shutdown(app2)
	test_utils.UpdateConfig(t, "user-mode", "true")
	session = test_utils.NewApiTestSession(t, app2)
	session.Post("/register", JSON{"name": "single", "email": "single@test.test", "password": "testpass"}, http.StatusOK)
	session.CheckResponse(nil)
	session.Get("/info", nil, http.StatusOK)
	body := session.Body()
	if Json(body)["team_id"] == nil {
		t.Fatal("Expected team_id")
	}
	test_utils.UpdateConfig(t, "user-mode", "false")

	// Email verification

	test_utils.UpdateConfig(t, "email-verification", "true")

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/register", JSON{"token": "AAA"}, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.InvalidJWT))

	session.Post("/register", JSON{"token": "AAA", "name": "test", "password": "testpass"}, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.InvalidJWT))

	key := []byte("test-secret-key")
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt_utils.Map{"a": "b"})
	signed, err := token.SignedString(key)
	if err != nil {
		t.Fatalf("failed to sign JWT: %v", err)
	}
	session.Post("/register", JSON{"token": signed, "name": "test", "password": "testpass"}, http.StatusUnauthorized)
	session.CheckResponse(errorf("token is unverifiable: error while executing keyfunc: " + consts.InvalidSigningAlgorithm))

	tokenStr, err := jwt_utils.GenerateJWT(t.Context(), jwt_utils.Map{"a": "b"})
	if err != nil {
		t.Fatalf("failed to generate JWT: %v", err)
	}
	session.Post("/register", JSON{"token": tokenStr, "name": "test", "password": "testpass"}, http.StatusUnauthorized)
	session.CheckResponse(errorf(consts.InvalidToken))

	tokenStr, err = jwt_utils.GenerateJWT(t.Context(), jwt_utils.Map{"email": "m@m.m"})
	if err != nil {
		t.Fatalf("failed to generate JWT: %v", err)
	}
	session.Post("/register", JSON{"token": tokenStr}, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.MissingRequiredFields))
	session.Post("/register", JSON{"token": tokenStr, "name": "m", "password": "12345678"}, http.StatusOK)
	session.CheckResponse(nil)

	expected := JSON{
		"email_verification": true,
		"name":               "m",
		"role":               "Player",
		"team_id":            nil,
		"user_mode":          false,
	}
	session.Get("/info", nil, http.StatusOK)
	session.CheckFilteredResponse(expected, "id")

	tokenStr, err = jwt_utils.GenerateJWT(t.Context(), jwt_utils.Map{"email": "m@m.m"})
	if err != nil {
		t.Fatalf("failed to generate JWT: %v", err)
	}
	session = test_utils.NewApiTestSession(t, app)
	session.Post("/register", JSON{"token": tokenStr, "name": "m2", "password": "testpass"}, http.StatusConflict)
	session.CheckResponse(errorf(consts.UserAlreadyExists))

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "m@m.m", "password": "12345678"}, http.StatusOK)
	session.CheckResponse(nil)
}
