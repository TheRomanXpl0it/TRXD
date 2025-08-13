package team_register_test

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
	test_utils.Main(m, "../../../", "team_register")
}

var testTeamRegister = []struct {
	testBody         interface{}
	secondUser       bool
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
		testBody:         JSON{"password": "testpass"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"name": "test", "password": strings.Repeat("a", consts.MinPasswordLength-1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.ShortPassword),
	},
	{
		testBody:         JSON{"name": "test", "password": strings.Repeat("a", consts.MaxPasswordLength+1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.LongPassword),
	},
	{
		testBody:         JSON{"name": strings.Repeat("a", consts.MaxNameLength+1), "password": "testpass"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.LongName),
	},
	{
		testBody:       JSON{"name": "test", "password": "testpass"},
		expectedStatus: http.StatusOK,
	},
	{
		testBody:         JSON{"name": "test", "password": "testpass"},
		expectedStatus:   http.StatusConflict,
		expectedResponse: errorf(consts.AlreadyInTeam),
	},
	{
		testBody:         JSON{"name": "test", "password": "testpass"},
		expectedStatus:   http.StatusConflict,
		secondUser:       true,
		expectedResponse: errorf(consts.TeamAlreadyExists),
	},
	{
		testBody:       JSON{"name": "test1", "password": "testpass"},
		expectedStatus: http.StatusOK,
		secondUser:     true,
	},
}

func TestTeamRegister(t *testing.T) {
	app := api.SetupApp()
	defer app.Shutdown()

	session := test_utils.NewApiTestSession(t, app)
	session.Post("/register", JSON{
		"username": "test",
		"email":    "test@test.test",
		"password": "testpass",
	}, http.StatusOK)
	session = test_utils.NewApiTestSession(t, app)
	session.Post("/register", JSON{
		"username": "test2",
		"email":    "test2@test.test",
		"password": "testpass",
	}, http.StatusOK)

	for _, test := range testTeamRegister {
		session := test_utils.NewApiTestSession(t, app)
		if test.secondUser {
			session.Post("/login", JSON{"email": "test2@test.test", "password": "testpass"}, http.StatusOK)
		} else {
			session.Post("/login", JSON{"email": "test@test.test", "password": "testpass"}, http.StatusOK)
		}
		session.Post("/teams", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}
}
