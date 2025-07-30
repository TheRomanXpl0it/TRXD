package tests

import (
	"net/http"
	"strings"
	"testing"
	"trxd/api"
	"trxd/db"
	"trxd/utils/consts"
)

var testRegister = []struct {
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
		testBody:         JSON{"username": "test"},
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
		testBody:         JSON{"username": "test", "email": "test@test.test"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"username": "test", "password": "testpass"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"email": "test@test.test", "password": "testpass"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"username": "test", "email": "test@test.test", "password": strings.Repeat("a", consts.MinPasswordLength-1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.ShortPassword),
	},
	{
		testBody:         JSON{"username": "test", "email": "test@test.test", "password": strings.Repeat("a", consts.MaxPasswordLength+1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.LongPassword),
	},
	{
		testBody:         JSON{"username": strings.Repeat("a", consts.MaxNameLength+1), "email": "test@test.test", "password": "testpass"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.LongName),
	},
	{
		testBody:         JSON{"username": "test", "email": strings.Repeat("a", consts.MaxEmailLength+1), "password": "testpass"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.LongEmail),
	},
	{
		testBody:         JSON{"username": "test", "email": "invalid-email", "password": "testpass"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidEmail),
	},
	{
		testBody:         JSON{"username": "test", "email": "test@test.test", "password": "testpass"},
		expectedStatus:   http.StatusOK,
		expectedResponse: JSON{"username": "test", "role": "Player"},
	},
	{
		testBody:         JSON{"username": "test", "email": "test@test.test", "password": "testpass"},
		expectedStatus:   http.StatusConflict,
		expectedResponse: errorf(consts.UserAlreadyExists),
	},
	{
		testBody:         JSON{"username": "test", "email": "test1@test.test", "password": "testpass"},
		expectedStatus:   http.StatusOK,
		expectedResponse: JSON{"username": "test", "role": "Player"},
	},
}

func TestRegister(t *testing.T) {
	db.DeleteAll()
	app := api.SetupApp()
	defer app.Shutdown()

	for _, test := range testRegister {
		session := newApiTestSession(t, app)
		session.Post("/register", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}
}
