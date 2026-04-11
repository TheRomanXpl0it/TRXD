package challenges_hidden_test

import (
	"math"
	"net/http"
	"testing"
	"trxd/api"
	"trxd/db/sqlc"
	"trxd/utils/consts"
	"trxd/utils/test_utils"
)

type JSON map[string]any

func errorf(val any) JSON {
	return JSON{"error": val}
}

func Json(val any) map[string]any {
	return val.(map[string]any)
}

func List(val any) []any {
	return val.([]any)
}

func Int32(val any) int32 {
	return int32(val.(float64))
}

func TestMain(m *testing.M) {
	test_utils.Main(m)
}

func TestRoute(t *testing.T) {
	app := api.SetupApp(t.Context())
	defer api.Shutdown(app)

	test_utils.RegisterUser(t, "author", "author@test.test", "authorpass", sqlc.UserRoleAuthor)

	session := test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "author@test.test", "password": "authorpass"}, http.StatusOK)
	session.Get("/challenges", nil, http.StatusOK)
	body := session.Body()

	challIDs := make([]int32, 0)
	hiddens := make(map[int32]bool)
	for _, chall := range List(body) {
		challIDs = append(challIDs, Int32(Json(chall)["id"]))
		hiddens[challIDs[len(challIDs)-1]] = Json(chall)["hidden"].(bool)
	}

	session.Patch("/challenges/hidden", JSON{}, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.MissingRequiredFields))

	session.Patch("/challenges/hidden", JSON{"chall_ids": []int32{}}, http.StatusOK)
	session.CheckResponse(nil)

	session.Patch("/challenges/hidden", JSON{"chall_ids": []int32{-1}}, http.StatusOK)
	session.CheckResponse(nil)

	session.Patch("/challenges/hidden", JSON{"chall_ids": []int32{99999}}, http.StatusOK)
	session.CheckResponse(nil)

	session.Patch("/challenges/hidden", JSON{"chall_ids": []int64{math.MaxInt32 + 1}}, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.InvalidJSON))

	session.Patch("/challenges/hidden", JSON{"chall_ids": challIDs}, http.StatusOK)
	session.CheckResponse(nil)

	session.Get("/challenges", nil, http.StatusOK)
	body = session.Body()

	for i, chall := range List(body) {
		hidden := hiddens[Int32(Json(chall)["id"])]
		if Json(chall)["hidden"].(bool) == hidden {
			t.Errorf("challenge with id %d should be hidden", challIDs[i])
		}
	}
}
