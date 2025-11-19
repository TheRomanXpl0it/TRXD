package categories_get_test

import (
	"net/http"
	"testing"
	"trxd/api"
	"trxd/utils/test_utils"
)

type JSON map[string]interface{}

func TestMain(m *testing.M) {
	test_utils.Main(m)
}

func TestRoute(t *testing.T) {
	app := api.SetupApp(t.Context())
	defer app.Shutdown()

	session := test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "admin@email.com", "password": "testpass"}, http.StatusOK)
	session.Post("/categories", JSON{"name": "cat-3", "icon": "icon-cat-3"}, http.StatusOK)

	expected := []JSON{
		{
			"icon": "cat-1",
			"name": "cat-1",
		},
		{
			"icon": "cat-2",
			"name": "cat-2",
		},
	}

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/register", JSON{"name": "test", "email": "test2@test.test", "password": "testpass"}, http.StatusOK)
	session.Post("/teams/register", JSON{"name": "test-team", "password": "testpass"}, http.StatusOK)
	session.Get("/categories", nil, http.StatusOK)
	session.CheckResponse(expected)

	expectedAdmin := []JSON{
		{
			"icon": "cat-1",
			"name": "cat-1",
		},
		{
			"icon": "cat-2",
			"name": "cat-2",
		},
		{
			"icon": "icon-cat-3",
			"name": "cat-3",
		},
	}

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "admin@email.com", "password": "testpass"}, http.StatusOK)
	session.Get("/categories", nil, http.StatusOK)
	session.CheckResponse(expectedAdmin)
}
