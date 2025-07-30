package tests

import (
	"net/http"
	"strings"
	"testing"
	"trxd/api"
	"trxd/db"
	"trxd/utils/consts"
)

var testJoinTeam = []struct {
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
		testBody:         JSON{"name": "test1", "password": "testpass"},
		expectedStatus:   http.StatusConflict,
		expectedResponse: errorf(consts.InvalidTeamCredentials),
	},
	{
		testBody:         JSON{"name": "test", "password": "testpassa"},
		expectedStatus:   http.StatusConflict,
		expectedResponse: errorf(consts.InvalidTeamCredentials),
	},
	{
		testBody:         JSON{"name": "test", "password": "testpass"},
		expectedStatus:   http.StatusOK,
		expectedResponse: JSON{"name": "test"},
	},
	{
		testBody:         JSON{"name": "test", "password": "testpass"},
		expectedStatus:   http.StatusConflict,
		secondUser:       true,
		expectedResponse: errorf(consts.AlreadyInTeam),
	},
}

func TestJoinTeam(t *testing.T) {
	db.DeleteAll()
	app := api.SetupApp()
	defer app.Shutdown()

	session := newApiTestSession(t, app)
	session.Post("/register", JSON{
		"email":    "test@test.test",
		"username": "test",
		"password": "testpass",
	}, http.StatusOK)
	session.Post("/register-team", JSON{
		"name":     "test",
		"password": "testpass",
	}, http.StatusOK)
	session.Post("/register", JSON{
		"email":    "test2@test.test",
		"username": "test2",
		"password": "testpass",
	}, http.StatusOK)

	for _, test := range testJoinTeam {
		session := newApiTestSession(t, app)
		if test.secondUser {
			session.Post("/login", JSON{"email": "test@test.test", "password": "testpass"}, http.StatusOK)
		} else {
			session.Post("/login", JSON{"email": "test2@test.test", "password": "testpass"}, http.StatusOK)
		}
		session.Post("/join-team", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}
}
