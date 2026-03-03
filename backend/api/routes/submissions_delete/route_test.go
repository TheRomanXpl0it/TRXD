package submissions_delete_test

import (
	"net/http"
	"testing"
	"trxd/api"
	"trxd/utils/test_utils"
)

//! TODO

type JSON map[string]interface{}

func errorf(val interface{}) JSON {
	return JSON{"error": val}
}

func TestMain(m *testing.M) {
	test_utils.Main(m)
}

func TestRoute(t *testing.T) {
	app := api.SetupApp(t.Context())
	defer api.Shutdown(app)

	Json := func(val interface{}) map[string]interface{} {
		return val.(map[string]interface{})
	}
	List := func(val interface{}) []interface{} {
		return val.([]interface{})
	}
	Int := func(val interface{}) int32 {
		return int32(val.(float64))
	} // TODO: make this more robust and global

	session := test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "admin@email.com", "password": "testpass"}, http.StatusOK)
	session.Get("/submissions", nil, http.StatusOK)
	body := session.Body() // TODO: make an optional argument (NoNull bool) (it throws an error if the body is nil) or make the opposite (Nullable bool) (so by default it is not nullable and throws an error if it is nil)
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	firstID := Int(Json(List(Json(body)["submissions"])[0])["id"])
	session.Delete("/submissions", JSON{"sub_id": firstID}, http.StatusOK)
	session.CheckResponse(nil)

	session.Get("/submissions", nil, http.StatusOK)
	body2 := session.Body()
	if body2 == nil {
		t.Fatal("Expected body to not be nil")
	}

	newFirstID := Int(Json(List(Json(body2)["submissions"])[0])["id"])
	if newFirstID >= firstID {
		t.Fatal("Expected first submission to be deleted: ", "firstID: ", firstID, " newFirstID: ", newFirstID)
	}

	len1 := len(List(Json(body)["submissions"]))
	len2 := len(List(Json(body2)["submissions"]))
	if len1 != len2+1 {
		t.Fatal("Expected number of submissions to decrease by 1: ", "before: ", len1, " after: ", len2)
	}

	submissions1 := List(Json(body)["submissions"])[1:]
	submissions2 := List(Json(body2)["submissions"])
	test_utils.Compare(t, submissions1, submissions2)
}
