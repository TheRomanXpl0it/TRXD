package configs_update_test

import (
	"net/http"
	"testing"
	"trxd/api"
	"trxd/db/sqlc"
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
		testBody:         JSON{"key": "allow-register"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"value": "true"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"key": "aaaaaaaaa", "value": "true"},
		expectedStatus:   http.StatusNotFound,
		expectedResponse: errorf(consts.ConfigNotFound),
	},
	{
		testBody:       JSON{"key": "allow-register", "value": "false"},
		expectedStatus: http.StatusOK,
	},
	{
		testBody:       JSON{"key": "allow-register", "value": "true"},
		expectedStatus: http.StatusOK,
	},
	{
		testBody:       JSON{"key": "allow-register", "value": "true"},
		expectedStatus: http.StatusOK,
	},
	{
		testBody:       JSON{"key": "domain", "value": ""},
		expectedStatus: http.StatusOK,
	},
	{
		testBody:       JSON{"key": "domain", "value": ""},
		expectedStatus: http.StatusOK,
	},
	{
		testBody:       JSON{"key": "domain", "value": "test.com"},
		expectedStatus: http.StatusOK,
	},
}

func TestRoute(t *testing.T) {
	app := api.SetupApp(t.Context())
	defer app.Shutdown()

	test_utils.RegisterUser(t, "test", "test@test.test", "testpass", sqlc.UserRoleAdmin)

	for _, test := range testData {
		session := test_utils.NewApiTestSession(t, app)
		session.Post("/login", JSON{"email": "test@test.test", "password": "testpass"}, http.StatusOK)
		session.Patch("/configs", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)

		if test.expectedStatus != http.StatusOK {
			continue
		}

		session.Get("/configs", nil, http.StatusOK)
		body := session.Body()
		if body == nil {
			t.Fatal("Expected response body to be non-nil")
		}
		for _, itemInt := range body.([]interface{}) {
			item := itemInt.(map[string]interface{})
			if item["key"] == test.testBody.(JSON)["key"] {
				if test.testBody.(JSON)["value"] != item["value"] {
					t.Fatalf("Values of config '%s', differs: %v != %v", item["key"], item["value"], test.testBody.(JSON)["value"])
				}
			}
		}
	}
}
