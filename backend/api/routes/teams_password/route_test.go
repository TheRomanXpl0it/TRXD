package teams_password_test

import (
	"fmt"
	"math"
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
	isAdmin          int
	teamNumber       int
	testBody         interface{}
	expectedStatus   int
	expectedResponse JSON
}{
	// Player && Team 0
	{
		isAdmin:          0,
		teamNumber:       0,
		testBody:         nil,
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidJSON),
	},
	{
		isAdmin:        0,
		teamNumber:     0,
		testBody:       JSON{},
		expectedStatus: http.StatusOK,
	},
	{
		isAdmin:          0,
		teamNumber:       0,
		testBody:         JSON{"team_id": -1},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(test_utils.Format(consts.MinError, "TeamID", 0)),
	},
	{
		isAdmin:          0,
		teamNumber:       0,
		testBody:         JSON{"team_id": math.MaxInt32 + 1},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidJSON),
	},
	{
		isAdmin:        0,
		teamNumber:     0,
		testBody:       JSON{"team_id": 0},
		expectedStatus: http.StatusOK,
	},
	{
		isAdmin:        0,
		teamNumber:     0,
		testBody:       JSON{"team_id": 0},
		expectedStatus: http.StatusOK,
	},
	{
		isAdmin:        0,
		teamNumber:     0,
		testBody:       JSON{"team_id": 0, "new_password": "NewPassw0rd!"},
		expectedStatus: http.StatusOK,
	},
	// Player && Team 1
	{
		isAdmin:          0,
		teamNumber:       1,
		testBody:         nil,
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidJSON),
	},
	{
		isAdmin:        0,
		teamNumber:     1,
		testBody:       JSON{},
		expectedStatus: http.StatusOK,
	},
	{
		isAdmin:          0,
		teamNumber:       1,
		testBody:         JSON{"team_id": -1},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(test_utils.Format(consts.MinError, "TeamID", 0)),
	},
	{
		isAdmin:          0,
		teamNumber:       1,
		testBody:         JSON{"team_id": math.MaxInt32 + 1},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidJSON),
	},
	{
		isAdmin:        0,
		teamNumber:     1,
		testBody:       JSON{"team_id": 0},
		expectedStatus: http.StatusOK,
	},
	{
		isAdmin:        0,
		teamNumber:     1,
		testBody:       JSON{"team_id": 0},
		expectedStatus: http.StatusOK,
	},
	{
		isAdmin:        0,
		teamNumber:     1,
		testBody:       JSON{"team_id": 0, "new_password": "NewPassw0rd!"},
		expectedStatus: http.StatusOK,
	},
	// Admin && Team 0
	{
		isAdmin:          1,
		teamNumber:       0,
		testBody:         nil,
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidJSON),
	},
	{
		isAdmin:          1,
		teamNumber:       0,
		testBody:         JSON{},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidTeamID),
	},
	{
		isAdmin:          1,
		teamNumber:       0,
		testBody:         JSON{"team_id": -1},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(test_utils.Format(consts.MinError, "TeamID", 0)),
	},
	{
		isAdmin:          1,
		teamNumber:       0,
		testBody:         JSON{"team_id": math.MaxInt32 + 1},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidJSON),
	},
	{
		isAdmin:        1,
		teamNumber:     0,
		testBody:       JSON{"team_id": 0},
		expectedStatus: http.StatusOK,
	},
	{
		isAdmin:        1,
		teamNumber:     0,
		testBody:       JSON{"team_id": 0},
		expectedStatus: http.StatusOK,
	},
	{
		isAdmin:        1,
		teamNumber:     0,
		testBody:       JSON{"team_id": 0, "new_password": "NewPassw0rd!"},
		expectedStatus: http.StatusOK,
	},
	// Admin && Team 1
	{
		isAdmin:          1,
		teamNumber:       1,
		testBody:         nil,
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidJSON),
	},
	{
		isAdmin:          1,
		teamNumber:       1,
		testBody:         JSON{},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidTeamID),
	},
	{
		isAdmin:          1,
		teamNumber:       1,
		testBody:         JSON{"team_id": -1},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(test_utils.Format(consts.MinError, "TeamID", 0)),
	},
	{
		isAdmin:          1,
		teamNumber:       1,
		testBody:         JSON{"team_id": math.MaxInt32 + 1},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidJSON),
	},
	{
		isAdmin:        1,
		teamNumber:     1,
		testBody:       JSON{"team_id": 0},
		expectedStatus: http.StatusOK,
	},
	{
		isAdmin:        1,
		teamNumber:     1,
		testBody:       JSON{"team_id": 0},
		expectedStatus: http.StatusOK,
	},
	{
		isAdmin:        1,
		teamNumber:     1,
		testBody:       JSON{"team_id": 0, "new_password": "NewPassw0rd!"},
		expectedStatus: http.StatusOK,
	},
}

func TestRoute(t *testing.T) {
	app := api.SetupApp(t.Context())
	defer api.Shutdown(app)

	user := test_utils.RegisterUser(t, "player", "player@test.test", "testpass", sqlc.UserRolePlayer)
	teamPlayer := test_utils.RegisterTeam(t, "playerTeam", "testpass", user.ID)
	test_utils.RegisterUser(t, "admin", "admin@test.test", "testpass", sqlc.UserRoleAdmin)
	user = test_utils.RegisterUser(t, "test", "test@test.test", "testpass", sqlc.UserRolePlayer)
	team := test_utils.RegisterTeam(t, "test", "testpass", user.ID)

	ids := []int32{teamPlayer.ID, team.ID}
	teams := []string{teamPlayer.Name, team.Name}
	passwords := []string{"testpass", "testpass"}

	for i, test := range testData {
		var email string
		if test.isAdmin == 1 {
			email = "admin@test.test"
		} else {
			email = "player@test.test"
		}

		session := test_utils.NewApiTestSession(t, app)
		session.Post("/login", JSON{"email": email, "password": "testpass"}, http.StatusOK)
		if body, ok := test.testBody.(JSON); ok && body != nil {
			if content, ok := body["team_id"]; ok && content == 0 {
				test.testBody.(JSON)["team_id"] = ids[test.teamNumber]
			}
		}
		session.Patch("/teams/password", test.testBody, test.expectedStatus)

		if test.isAdmin == 0 && test.teamNumber == 1 {
			test.teamNumber = 0
		}

		var newPass interface{}
		var passOk bool
		if test.testBody != nil {
			newPass, passOk = test.testBody.(JSON)["new_password"]
		}
		if test.expectedStatus != http.StatusOK || passOk {
			session.CheckResponse(test.expectedResponse)
		}
		if test.expectedStatus == http.StatusOK {
			if passOk {
				passwords[test.teamNumber] = newPass.(string)
			} else {
				sessionBody := session.Body()
				if sessionBody == nil {
					t.Fatal("Expected body to not be nil")
				}
				body := sessionBody.(map[string]interface{})
				newPasswordInterface, ok := body["new_password"]
				if !ok {
					t.Fatalf("Expected 'new_password' in response, got: %v", body)
				}
				passwords[test.teamNumber], ok = newPasswordInterface.(string)
				if !ok {
					t.Fatalf("Expected 'new_password' to be a string, got: %T", newPasswordInterface)
				}
			}
		}

		idx := i + (test.teamNumber * len(testData) / 4) + (test.isAdmin * len(testData) / 2)
		session = test_utils.NewApiTestSession(t, app)
		session.Post("/register", JSON{"name": fmt.Sprintf("test%d", idx), "email": fmt.Sprintf("test%d@test.test", idx), "password": "testpass"}, http.StatusOK)
		session.Post("/teams/join", JSON{"name": teams[test.teamNumber], "password": passwords[test.teamNumber]}, http.StatusOK)
		session.CheckResponse(nil)
		session.Post("/submissions", JSON{}, http.StatusBadRequest)
	}
}
