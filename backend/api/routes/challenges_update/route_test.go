package challenges_update_test

import (
	"fmt"
	"net/http"
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
	testFiles        []string
	expectedStatus   int
	expectedResponse JSON
}{
	{
		testBody:         nil,
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.ChallIDRequired),
	},
	{
		testBody:         JSON{"chall_id": ""},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.NoDataToUpdate),
	},
	{
		testBody:         JSON{"chall_id": "", "name": strings.Repeat("a", consts.MaxChallNameLength+1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.LongChallName),
	},
	{
		testBody:         JSON{"chall_id": "", "category": strings.Repeat("a", consts.MaxCategoryLength+1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.LongCategory),
	},
	{
		testBody:         JSON{"chall_id": "", "description": strings.Repeat("a", consts.MaxChallDescLength+1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.LongChallDesc),
	},
	{
		testBody:         JSON{"chall_id": "", "difficulty": strings.Repeat("a", consts.MaxChallDifficultyLength+1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.LongChallDifficulty),
	},
	{
		testBody:         JSON{"chall_id": "", "max_points": -1},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidChallMaxPoints),
	},
	{
		testBody:         JSON{"chall_id": "", "port": consts.MinPort - 1},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidPort),
	},
	{
		testBody:         JSON{"chall_id": "", "port": consts.MaxPort + 1},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidPort),
	},
	{
		testBody:         JSON{"chall_id": "", "lifetime": 0},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidLifetime),
	},
	{
		testBody:         JSON{"chall_id": "", "envs": "<invalid-json>"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidEnvs),
	},
	{
		testBody:         JSON{"chall_id": "", "max_memory": 0},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidMaxMemory),
	},
	{
		testBody:         JSON{"chall_id": "", "max_cpu": "<invalid-float>"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidMaxCpu),
	},
	{
		testBody:         JSON{"chall_id": -1, "name": "test"},
		expectedStatus:   http.StatusNotFound,
		expectedResponse: errorf(consts.ChallengeNotFound),
	},
	{
		testBody:         JSON{"chall_id": "", "name": "chall-2"},
		expectedStatus:   http.StatusConflict,
		expectedResponse: errorf(consts.ChallNameExists),
	},
	{
		testBody:         JSON{"chall_id": "", "category": "<invalid-category>"},
		expectedStatus:   http.StatusNotFound,
		expectedResponse: errorf(consts.CategoryNotFound),
	},
	{
		testBody: JSON{
			"chall_id":    "",
			"name":        "Test",
			"category":    "cat-2",
			"description": "new test desc",
			"difficulty":  "Ez",
			"authors":     []string{"author1", "author2"},
			"type":        "Container",
			"hidden":      false,
			"max_points":  1000,
			"score_type":  "Dynamic",
			"host":        "http://ctf.theromanxpl0.it",
			"port":        1234,

			"image":       "ubuntu:latest",
			"compose":     "",
			"hash_domain": true,
			"lifetime":    60,
			"envs":        `{"key": "value"}`,
			"max_memory":  512,
			"max_cpu":     "1.0",
		},
		testFiles:      []string{"f1.txt", "f2.txt", "f3.txt"},
		expectedStatus: http.StatusOK,
	},
}

func TestRoute(t *testing.T) {
	app := api.SetupApp()
	defer app.Shutdown()

	module := test_utils.GetModuleName(t)
	dir := "/tmp/" + module + "/"
	test_utils.CreateDir(t, dir)
	test_utils.CreateFile(t, dir+"f1.txt", "test-line 1\ntest-line 2")
	test_utils.CreateFile(t, dir+"f2.txt", "")
	test_utils.CreateFile(t, dir+"f3.txt", string([]byte{0, 0, 255, 255, 127, 97}))

	test_utils.RegisterUser(t, "author", "author@test.test", "authorpass", sqlc.UserRoleAuthor)

	var challID int32
	session := test_utils.NewApiTestSession(t, app)
	session.Post("/users/login", JSON{"email": "author@test.test", "password": "authorpass"}, http.StatusOK)
	session.Post("/categories/create", JSON{"name": "cat", "icon": "icon"}, -1)

	chall := test_utils.TryCreateChallenge(t, "chall", "cat", "test-desc", sqlc.DeployTypeNormal, 1, sqlc.ScoreTypeStatic)
	if chall != nil {
		challID = chall.ID
	}

	session.Patch("/challenges/update", nil, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.InvalidMultipartForm))

	session.RequestMultipart(http.MethodPatch, "/challenges/update", JSON{"chall_id": challID}, []string{dir + "f1.txt"}, []string{"/"}, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.InvalidFilePath))

	for i, test := range testData {
		session := test_utils.NewApiTestSession(t, app)
		session.Post("/users/login", JSON{"email": "author@test.test", "password": "authorpass"}, http.StatusOK)

		if test.testBody != nil {
			if content, ok := test.testBody["chall_id"]; ok && content == "" {
				test.testBody["chall_id"] = challID
			}
		}

		for i := range test.testFiles {
			test.testFiles[i] = dir + test.testFiles[i]
		}

		session.PatchMultipart("/challenges/update", test.testBody, test.testFiles, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)

		if i == len(testData)-1 {
			session.Get(fmt.Sprintf("/challenges/%d", challID), nil, http.StatusOK)
			body := session.Body()
			expected := JSON{
				"attachments": []string{
					fmt.Sprintf("attachments/%d/%s", challID, "f1.txt"),
					fmt.Sprintf("attachments/%d/%s", challID, "f2.txt"),
					fmt.Sprintf("attachments/%d/%s", challID, "f3.txt"),
				},
				"authors":     test.testBody["authors"],
				"category":    test.testBody["category"],
				"description": test.testBody["description"],
				"difficulty":  test.testBody["difficulty"],
				"docker_config": JSON{
					"envs":        test.testBody["envs"],
					"hash_domain": test.testBody["hash_domain"],
					"image":       test.testBody["image"],
					"lifetime":    test.testBody["lifetime"],
					"max_cpu":     test.testBody["max_cpu"],
					"max_memory":  test.testBody["max_memory"],
				},
				"first_blood": nil,
				"flags":       []string{},
				"hidden":      test.testBody["hidden"],
				"host":        test.testBody["host"],
				"id":          challID,
				"instance":    test.testBody["type"] != "Normal",
				"max_points":  test.testBody["max_points"],
				"name":        test.testBody["name"],
				"points":      test.testBody["max_points"],
				"port":        test.testBody["port"],
				"score_type":  test.testBody["score_type"],
				"solved":      false,
				"solves":      0,
				"solves_list": []string{},
				"tags":        []string{},
				"timeout":     0,
				"type":        test.testBody["type"],
			}
			test_utils.Compare(t, expected, body)
		}
	}

}
