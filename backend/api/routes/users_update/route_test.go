package users_update_test

import (
	"fmt"
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
		testBody:         JSON{},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"name": strings.Repeat("a", consts.MaxUserNameLen+1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(test_utils.Format(consts.MaxError, "Name", consts.MaxUserNameLen)),
	},
	{
		testBody:         JSON{"country": "a"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidCountry),
	},
	{
		testBody:         JSON{"name": "a", "country": "USA"},
		expectedStatus:   http.StatusConflict,
		expectedResponse: errorf(consts.NameAlreadyTaken),
	},
	{
		testBody:       JSON{"name": "aa", "country": "USA"},
		expectedStatus: http.StatusOK,
	},
	{
		testBody:       JSON{"name": "aa", "country": "USA"},
		expectedStatus: http.StatusOK,
	},
	{
		testBody:       JSON{"name": "aa", "country": ""},
		expectedStatus: http.StatusOK,
	},
	{
		testBody:       JSON{"name": "bb", "country": "JPN"},
		expectedStatus: http.StatusOK,
	},
	{
		testBody:       JSON{"name": "cc"},
		expectedStatus: http.StatusOK,
	},
}

func TestRoute(t *testing.T) {
	app := api.SetupApp(t.Context())
	defer api.Shutdown(app)

	session := test_utils.NewApiTestSession(t, app)
	session.Post("/register", JSON{"name": "test", "email": "test@test.test", "password": "testpass"}, http.StatusOK)
	session.CheckResponse(nil)

	for _, test := range testData {
		session := test_utils.NewApiTestSession(t, app)
		session.Post("/login", JSON{"email": "test@test.test", "password": "testpass"}, http.StatusOK)
		session.Patch("/users", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}

	app2 := api.SetupApp(t.Context())
	defer api.Shutdown(app2)
	test_utils.UpdateConfig(t, "user-mode", "true")

	session = test_utils.NewApiTestSession(t, app2)
	session.Post("/login", JSON{"email": "test@test.test", "password": "testpass"}, http.StatusOK)
	session.CheckResponse(nil)
	session.Get("/info", nil, http.StatusOK)
	body := session.Body()
	if body == nil {
		t.Fatal("Expected body")
	}
	if body.(map[string]interface{})["team_id"] != nil {
		t.Fatal("Expected no team_id")
	}

	session.Post("/teams/register", JSON{"name": "team", "password": "teampass"}, http.StatusOK)
	session.CheckResponse(nil)
	session.Patch("/users", JSON{"name": "updated-name", "country": "USA"}, http.StatusOK)
	session.CheckResponse(nil)

	session.Get("/info", nil, http.StatusOK)
	body = session.Body()
	if body == nil {
		t.Fatal("Expected body")
	}
	tid := body.(map[string]interface{})["team_id"]
	if tid == nil {
		t.Fatal("Expected team_id")
	}

	session.Get(fmt.Sprintf("/teams/%v", int(tid.(float64))), nil, http.StatusOK)
	body = session.Body()
	if body == nil {
		t.Fatal("Expected body")
	}
	if body.(map[string]interface{})["name"] != "updated-name" {
		t.Fatal("Expected updated name")
	}
	if body.(map[string]interface{})["country"] != "USA" {
		t.Fatal("Expected updated country")
	}
}
