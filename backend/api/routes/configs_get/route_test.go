package configs_get_test

import (
	"net/http"
	"testing"
	"trxd/api"
	"trxd/db/sqlc"
	"trxd/utils/test_utils"
)

type JSON map[string]interface{}

func TestMain(m *testing.M) {
	test_utils.Main(m)
}

func TestRoute(t *testing.T) {
	app := api.SetupApp()
	defer app.Shutdown()

	test_utils.RegisterUser(t, "test", "test@test.test", "testpass", sqlc.UserRoleAdmin)
	session := test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "test@test.test", "password": "testpass"}, http.StatusOK)
	session.Get("/configs", nil, http.StatusOK)
	session.CheckResponse([]JSON{
		{
			"description": "",
			"key":         "allow-register",
			"type":        "bool",
			"value":       "true",
		},
		{
			"description": "",
			"key":         "chall-min-points",
			"type":        "int",
			"value":       "50",
		},
		{
			"description": "",
			"key":         "chall-points-decay",
			"type":        "int",
			"value":       "15",
		},
		{
			"description": "",
			"key":         "domain",
			"type":        "string",
			"value":       "",
		},
		{
			"description": "",
			"key":         "hash-len",
			"type":        "int",
			"value":       "12",
		},
		{
			"description": "",
			"key":         "instance-lifetime",
			"type":        "int",
			"value":       "1800",
		},
		{
			"description": "",
			"key":         "instance-max-cpu",
			"type":        "string",
			"value":       "1.0",
		},
		{
			"description": "",
			"key":         "instance-max-mem",
			"type":        "int",
			"value":       "512",
		},
		{
			"description": "",
			"key":         "max-port",
			"type":        "int",
			"value":       "30000",
		},
		{
			"description": "",
			"key":         "min-port",
			"type":        "int",
			"value":       "20000",
		},
		{
			"description": "",
			"key":         "secret",
			"type":        "string",
			"value":       "",
		},
	})
}
