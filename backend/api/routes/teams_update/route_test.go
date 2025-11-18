package teams_update_test

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
		testBody:         JSON{},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"country": "a"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidCountry),
	},
	{
		testBody:         JSON{"image": "a"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidHttpUrl),
	},
	{
		testBody:         JSON{"image": "file://example.com/image.png"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidHttpUrl),
	},
	{
		testBody:         JSON{"bio": strings.Repeat("a", consts.MaxBioLen+1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf("Bio must not exceed 10240"),
	},
	{
		testBody:         JSON{"name": strings.Repeat("a", consts.MaxTeamNameLen+1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf("Name must not exceed 64"),
	},
	{
		testBody:         JSON{"name": "A"},
		expectedStatus:   http.StatusConflict,
		expectedResponse: errorf(consts.NameAlreadyTaken),
	},
	{
		testBody:       JSON{"country": "USA", "image": "http://example.com/image.png", "bio": "a"},
		expectedStatus: http.StatusOK,
	},
	{
		testBody:       JSON{"country": "USA", "image": "https://example.com/image.png", "bio": "a"},
		expectedStatus: http.StatusOK,
	},
	{
		testBody:       JSON{"country": "", "image": "", "bio": ""},
		expectedStatus: http.StatusOK,
	},
	{
		testBody:       JSON{"country": "JPN", "image": "https://example.com/image.png"},
		expectedStatus: http.StatusOK,
	},
	{
		testBody:       JSON{"country": "ITA"},
		expectedStatus: http.StatusOK,
	},
}

func TestRoute(t *testing.T) {
	app := api.SetupApp(t.Context())
	defer app.Shutdown()

	session := test_utils.NewApiTestSession(t, app)
	session.Post("/register", JSON{"name": "test", "email": "test@test.test", "password": "testpass"}, http.StatusOK)
	session.Post("/teams/register", JSON{"name": "test-team", "password": "testpass"}, http.StatusOK)

	for _, test := range testData {
		session := test_utils.NewApiTestSession(t, app)
		session.Post("/login", JSON{"email": "test@test.test", "password": "testpass"}, http.StatusOK)
		session.Patch("/teams", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}
}
