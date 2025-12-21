package users_register_test

import (
	"net/http"
	"strings"
	"testing"
	"trxd/api"
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

	app2 := api.SetupApp(t.Context())
	defer api.Shutdown(app2)
	test_utils.UpdateConfig(t, "user-mode", "true")
	session = test_utils.NewApiTestSession(t, app2)
	session.Post("/register", JSON{"name": "single", "email": "single@test.test", "password": "testpass"}, http.StatusOK)
	session.CheckResponse(nil)
	session.Get("/info", nil, http.StatusOK)
	body := session.Body()
	if body == nil {
		t.Fatal("Expected body")
	}
	if body.(map[string]interface{})["team_id"] == nil {
		t.Fatal("Expected team_id")
	}
}
