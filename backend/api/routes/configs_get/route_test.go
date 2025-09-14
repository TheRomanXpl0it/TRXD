package configs_get_test

import (
	"fmt"
	"net/http"
	"sort"
	"testing"
	"trxd/api"
	"trxd/db/sqlc"
	"trxd/utils/consts"
	"trxd/utils/test_utils"
)

type JSON map[string]interface{}

func TestMain(m *testing.M) {
	test_utils.Main(m)
}

func TestRoute(t *testing.T) {
	app := api.SetupApp()
	defer app.Shutdown()

	expected := make([]JSON, 0, len(consts.DefaultConfigs))
	for key, value := range consts.DefaultConfigs {
		expected = append(expected, JSON{
			"description": "",
			"key":         key,
			"type":        fmt.Sprintf("%T", value),
			"value":       fmt.Sprintf("%v", value),
		})
	}

	sort.Slice(expected, func(i int, j int) bool {
		return expected[i]["key"].(string) < expected[j]["key"].(string)
	})

	test_utils.RegisterUser(t, "test", "test@test.test", "testpass", sqlc.UserRoleAdmin)
	session := test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "test@test.test", "password": "testpass"}, http.StatusOK)
	session.Get("/configs", nil, http.StatusOK)
	session.CheckResponse(expected)
}
