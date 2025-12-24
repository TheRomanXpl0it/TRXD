package attachments_create_test

import (
	"fmt"
	"math"
	"net/http"
	"os"
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
	testBody         JSON
	testFiles        []string
	expectedStatus   int
	expectedResponse JSON
}{
	{
		testBody:         nil,
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{},
		testFiles:        []string{"f1.txt", "f2.txt", "f3.txt"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"chall_id": ""},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"chall_id": -1},
		testFiles:        []string{"f1.txt", "f2.txt", "f3.txt"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(test_utils.Format(consts.MinError, "ChallID", 0)),
	},
	{
		testBody:         JSON{"chall_id": 99999},
		testFiles:        []string{"f1.txt", "f2.txt", "f3.txt"},
		expectedStatus:   http.StatusNotFound,
		expectedResponse: errorf(consts.ChallengeNotFound),
	},
	{
		testBody:         JSON{"chall_id": math.MaxInt32 + 1},
		testFiles:        []string{"f1.txt", "f2.txt", "f3.txt"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidFormData),
	},
	{
		testBody:       JSON{"chall_id": ""},
		testFiles:      []string{"f1.txt", "f2.txt", "f3.txt"},
		expectedStatus: http.StatusOK,
	},
}

func TestRoute(t *testing.T) {
	app := api.SetupApp(t.Context())
	defer api.Shutdown(app)

	module := test_utils.GetModuleName(t)
	dir := "/tmp/" + module + "/"
	test_utils.CreateDir(t, dir)
	test_utils.CreateFile(t, dir+"f1.txt", "test-line 1\ntest-line 2")
	test_utils.CreateFile(t, dir+"f2.txt", "")
	test_utils.CreateFile(t, dir+"f3.txt", string([]byte{0, 0, 255, 255, 127, 97}))

	h1 := test_utils.HashFile(t, dir+"f1.txt")
	h2 := test_utils.HashFile(t, dir+"f2.txt")
	h3 := test_utils.HashFile(t, dir+"f3.txt")

	test_utils.RegisterUser(t, "author", "author@test.test", "authorpass", sqlc.UserRoleAuthor)

	var challID int32
	session := test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "author@test.test", "password": "authorpass"}, http.StatusOK)
	session.Post("/categories", JSON{"name": "cat"}, -1)

	chall := test_utils.TryCreateChallenge(t, "chall", "cat", "test-desc", sqlc.DeployTypeNormal, 1, sqlc.ScoreTypeStatic)
	if chall != nil {
		challID = chall.ID
	}

	session.Post("/attachments", nil, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.InvalidMultipartForm))

	session.RequestMultipart(http.MethodPost, "/attachments", JSON{"chall_id": challID}, []string{dir + "f1.txt"}, []string{"/"}, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.InvalidFilePath))

	for i, test := range testData {
		session := test_utils.NewApiTestSession(t, app)
		session.Post("/login", JSON{"email": "author@test.test", "password": "authorpass"}, http.StatusOK)

		if test.testBody != nil {
			if content, ok := test.testBody["chall_id"]; ok && content == "" {
				test.testBody["chall_id"] = challID
			}
		}

		for i := range test.testFiles {
			test.testFiles[i] = dir + test.testFiles[i]
		}

		session.PostMultipart("/attachments", test.testBody, test.testFiles, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)

		if i != len(testData)-1 {
			continue
		}

		session.Get("/challenges", nil, http.StatusOK)
		body := session.Body()
		if body == nil {
			t.Fatal("Expected body to not be nil")
		}
		expected := JSON{
			"attachments": []string{
				fmt.Sprintf("/%s/%s", h1, "f1.txt"),
				fmt.Sprintf("/%s/%s", h2, "f2.txt"),
				fmt.Sprintf("/%s/%s", h3, "f3.txt"),
			},
			"authors":     []string{},
			"category":    "cat",
			"conn_type":   "NONE",
			"description": "test-desc",
			"difficulty":  "",
			"first_blood": false,
			"hidden":      true,
			"host":        "",
			"id":          challID,
			"instance":    false,
			"max_points":  1,
			"name":        "chall",
			"points":      1,
			"port":        0,
			"score_type":  sqlc.ScoreTypeStatic,
			"solved":      false,
			"solves":      0,
			"tags":        []string{},
			"timeout":     0,
		}
		var challengeBody interface{}
		for _, v := range body.([]interface{}) {
			if int32(v.(map[string]interface{})["id"].(float64)) == challID {
				challengeBody = v
				break
			}
		}
		test_utils.Compare(t, expected, challengeBody)

		attachments := expected["attachments"].([]string)
		for _, name := range attachments {
			path := fmt.Sprintf("attachments/%d/%s", challID, name)
			if _, err := os.Stat(path); os.IsNotExist(err) {
				t.Fatalf("Expected attachment file %s to exist", path)
			}
		}
	}
}
