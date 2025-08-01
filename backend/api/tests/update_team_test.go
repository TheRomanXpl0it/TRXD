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

var testUpdateTeam = []struct {
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
		testBody:         JSON{},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"nationality": strings.Repeat("a", consts.MaxNationalityLength+1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.LongNationality),
	},
	{
		testBody:       JSON{"nationality": "a", "image": "a", "bio": "a"},
		expectedStatus: http.StatusOK,
	},
	{
		testBody:       JSON{"nationality": "b", "image": "b"},
		expectedStatus: http.StatusOK,
	},
	{
		testBody:       JSON{"nationality": "c"},
		expectedStatus: http.StatusOK,
	},
}

func TestUpdateTeam(t *testing.T) {
	db.DeleteAll()
	app := api.SetupApp()
	defer app.Shutdown()

	session := utils.NewApiTestSession(t, app)
	session.Post("/register", JSON{"username": "test", "email": "test@test.test", "password": "testpass"}, http.StatusOK)
	session.Post("/register-team", JSON{"name": "test-team", "password": "testpass"}, http.StatusOK)

	for _, test := range testUpdateTeam {
		session := utils.NewApiTestSession(t, app)
		session.Post("/login", JSON{"email": "test@test.test", "password": "testpass"}, http.StatusOK)
		session.Patch("/team", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}
}
