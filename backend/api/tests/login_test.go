package tests

import (
	"net/http"
	"strings"
	"testing"
	"trxd/api"
	"trxd/db"
	"trxd/utils"
	"trxd/utils/consts"
)

var testLogin = []struct {
	testBody         interface{}
	register         bool
	expectedStatus   int
	expectedResponse JSON
}{
	{
		testBody:         nil,
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidJSON),
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
		testBody:         JSON{"email": "test@test.test", "password": strings.Repeat("a", consts.MaxPasswordLength+1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.LongPassword),
	},
	{
		testBody:         JSON{"email": strings.Repeat("a", consts.MaxEmailLength+1), "password": "testpass"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.LongEmail),
	},
	{
		testBody:         JSON{"email": "test@test.test", "password": "testpass"},
		expectedStatus:   http.StatusUnauthorized,
		expectedResponse: errorf(consts.InvalidCredentials),
	},
	{
		testBody:         JSON{"username": "test", "email": "test@test.test", "password": "testpass"},
		register:         true,
		expectedStatus:   http.StatusOK,
		expectedResponse: JSON{"username": "test", "role": string(db.UserRolePlayer)},
	},
}

func TestLogin(t *testing.T) {
	db.DeleteAll()
	app := api.SetupApp()
	defer app.Shutdown()

	for _, test := range testLogin {

		if test.register {
			session := utils.NewApiTestSession(t, app)
			session.Post("/register", test.testBody, http.StatusOK)
		}

		session := utils.NewApiTestSession(t, app)
		session.Post("/login", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}
}
