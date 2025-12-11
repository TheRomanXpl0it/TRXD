package attachments_delete_test

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"strings"
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
	expectedStatus   int
	expectedResponse JSON
}{
	{
		testBody:         JSON{},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"chall_id": ""},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"names": []string{"test"}},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"chall_id": "", "names": []string{strings.Repeat("a", consts.MaxAttachmentNameLen+1)}},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(test_utils.Format(consts.MaxError, "Names[0]", consts.MaxAttachmentNameLen)),
	},
	{
		testBody:         JSON{"chall_id": -1, "names": []string{"test"}},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(test_utils.Format(consts.MinError, "ChallID", 0)),
	},
	{
		testBody:         JSON{"chall_id": 99999, "names": []string{"test"}},
		expectedStatus:   http.StatusNotFound,
		expectedResponse: errorf(consts.AttachmentNotFound),
	},
	{
		testBody:         JSON{"chall_id": math.MaxInt32 + 1, "names": []string{"test"}},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidJSON),
	},
	{
		testBody:         JSON{"chall_id": "", "names": []string{"aaaaaaaaaaaaa"}},
		expectedStatus:   http.StatusNotFound,
		expectedResponse: errorf(consts.AttachmentNotFound),
	},
	{
		testBody:       JSON{"chall_id": "", "names": []string{"test1.txt"}},
		expectedStatus: http.StatusOK,
	},
	{
		testBody:         JSON{"chall_id": "", "names": []string{"test1.txt"}},
		expectedStatus:   http.StatusNotFound,
		expectedResponse: errorf(consts.AttachmentNotFound),
	},
	{
		testBody:       JSON{"chall_id": "", "names": []string{"test2.txt"}},
		expectedStatus: http.StatusOK,
	},
}

func TestRoute(t *testing.T) {
	app := api.SetupApp(t.Context())
	defer app.Shutdown()

	module := test_utils.GetModuleName(t)
	dir := "/tmp/" + module + "/"
	test_utils.CreateDir(t, dir)
	test_utils.CreateFile(t, dir+"test1.txt", "test-line 1\ntest-line 2")
	test_utils.CreateFile(t, dir+"test2.txt", "")
	test_utils.CreateFile(t, dir+"test3.txt", string([]byte{0, 0, 255, 255, 127, 97}))

	test_utils.RegisterUser(t, "author", "author@test.test", "authorpass", sqlc.UserRoleAuthor)

	var challID int32
	session := test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "author@test.test", "password": "authorpass"}, http.StatusOK)
	session.Post("/categories", JSON{"name": "cat"}, -1)

	chall := test_utils.TryCreateChallenge(t, "chall", "cat", "test-desc", sqlc.DeployTypeNormal, 1, sqlc.ScoreTypeStatic)
	if chall != nil {
		challID = chall.ID
	}

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "author@test.test", "password": "authorpass"}, http.StatusOK)
	session.PostMultipart("/attachments", JSON{"chall_id": challID}, []string{dir + "test1.txt", dir + "test2.txt", dir + "test3.txt"}, http.StatusOK)
	session.CheckResponse(nil)

	for i, test := range testData {
		session := test_utils.NewApiTestSession(t, app)
		session.Post("/login", JSON{"email": "author@test.test", "password": "authorpass"}, http.StatusOK)

		if test.testBody != nil {
			if content, ok := test.testBody["chall_id"]; ok && content == "" {
				test.testBody["chall_id"] = challID
			}
		}

		session.Delete("/attachments", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)

		if test.expectedStatus != http.StatusOK {
			continue
		}

		var attachments []string
		if i == len(testData)-1 {
			attachments = []string{
				fmt.Sprintf("/%s/%s", "test3.txt", "test3.txt"),
			}
		} else {
			attachments = []string{
				fmt.Sprintf("/%s/%s", "test2.txt", "test2.txt"),
				fmt.Sprintf("/%s/%s", "test3.txt", "test3.txt"),
			}
		}

		session.Get("/challenges", nil, http.StatusOK)
		body := session.Body()
		if body == nil {
			t.Fatal("Expected body to not be nil")
		}
		expected := JSON{
			"attachments": attachments,
			"authors":     []string{},
			"category":    "cat",
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

		path := fmt.Sprintf("attachments/%d/%s", challID, test.testBody["names"].([]string)[0])
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Fatalf("Expected attachment file %s to be deleted", path)
		}
		for _, name := range attachments {
			path := fmt.Sprintf("attachments/%d/%s", challID, name)
			if _, err := os.Stat(path); os.IsNotExist(err) {
				t.Fatalf("Expected attachment file %s to exist", path)
			}
		}
	}
}
