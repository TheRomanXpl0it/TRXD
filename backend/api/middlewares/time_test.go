package middlewares_test

import (
	"net/http"
	"testing"
	"time"
	"trxd/api"
	"trxd/db/sqlc"
	"trxd/utils/consts"
	"trxd/utils/test_utils"
)

func TestTimeMiddlewares(t *testing.T) {
	app := api.SetupApp(t.Context())
	defer api.Shutdown(app)

	test_utils.RegisterUser(t, "ADMIN", "admin@admin.com", "testpass", sqlc.UserRoleAdmin)
	test_utils.RegisterUser(t, "PLAYER", "player@player.com", "testpass", sqlc.UserRolePlayer)

	admin := test_utils.NewApiTestSession(t, app)
	admin.Post("/login", JSON{"email": "admin@admin.com", "password": "testpass"}, http.StatusOK)
	player := test_utils.NewApiTestSession(t, app)
	player.Post("/login", JSON{"email": "player@player.com", "password": "testpass"}, http.StatusOK)
	player.Post("/teams/register", JSON{"name": "team-player", "password": "teampass"}, http.StatusOK)
	player.CheckResponse(nil)

	admin.Get("/challenges", nil, http.StatusOK)
	admin.Get("/challenges/99999", nil, http.StatusNotFound)
	admin.Post("/submissions", nil, http.StatusBadRequest)
	player.Get("/challenges", nil, http.StatusOK)
	player.Get("/challenges/99999", nil, http.StatusNotFound)
	player.Post("/submissions", nil, http.StatusBadRequest)

	startTime := time.Now().Add(10 * time.Second).Format(time.RFC3339)
	admin.Patch("/configs", JSON{"key": "start-time", "value": startTime}, http.StatusOK)
	admin.CheckResponse(nil)

	admin.Get("/challenges", nil, http.StatusOK)
	admin.Get("/challenges/99999", nil, http.StatusNotFound)
	admin.Post("/submissions", nil, http.StatusBadRequest)
	player.Get("/challenges", nil, http.StatusForbidden)
	player.CheckResponse(errorf(consts.NotStartedYet))
	player.Get("/challenges/99999", nil, http.StatusForbidden)
	player.CheckResponse(errorf(consts.NotStartedYet))
	player.Post("/submissions", nil, http.StatusForbidden)
	player.CheckResponse(errorf(consts.NotStartedYet))

	time.Sleep(12 * time.Second)

	admin.Get("/challenges", nil, http.StatusOK)
	admin.Get("/challenges/99999", nil, http.StatusNotFound)
	admin.Post("/submissions", nil, http.StatusBadRequest)
	player.Get("/challenges", nil, http.StatusOK)
	player.Get("/challenges/99999", nil, http.StatusNotFound)
	player.Post("/submissions", nil, http.StatusBadRequest)

	endTime := time.Now().Add(-10 * time.Second).Format(time.RFC3339)
	admin.Patch("/configs", JSON{"key": "end-time", "value": endTime}, http.StatusOK)
	admin.CheckResponse(nil)

	admin.Post("/submissions", nil, http.StatusBadRequest)
	player.Post("/submissions", nil, http.StatusForbidden)
	player.CheckResponse(errorf(consts.AlreadyEnded))
}
