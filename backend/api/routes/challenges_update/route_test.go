package challenges_update_test

import (
	"fmt"
	"math"
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
	expectedStatus   int
	expectedResponse JSON
}{
	{
		testBody:         nil,
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"chall_id": ""},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.NoDataToUpdate),
	},
	{
		testBody:         JSON{"chall_id": "", "name": strings.Repeat("a", consts.MaxChallNameLen+1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(test_utils.Format(consts.MaxError, "Name", consts.MaxChallNameLen)),
	},
	{
		testBody:         JSON{"chall_id": "", "category": strings.Repeat("a", consts.MaxCategoryLen+1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(test_utils.Format(consts.MaxError, "Category", consts.MaxCategoryLen)),
	},
	{
		testBody:         JSON{"chall_id": "", "description": strings.Repeat("a", consts.MaxChallDescLen+1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(test_utils.Format(consts.MaxError, "Description", consts.MaxChallDescLen)),
	},
	{
		testBody:         JSON{"chall_id": "", "difficulty": strings.Repeat("a", consts.MaxChallDifficultyLen+1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(test_utils.Format(consts.MaxError, "Difficulty", consts.MaxChallDifficultyLen)),
	},
	{
		testBody:         JSON{"chall_id": "", "authors": ""},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidJSON),
	},
	{
		testBody:         JSON{"chall_id": "", "authors": []string{strings.Repeat("a", consts.MaxAuthorNameLen+1)}},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(test_utils.Format(consts.MaxError, "Authors[0]", consts.MaxAuthorNameLen)),
	},
	{
		testBody:         JSON{"chall_id": "", "tags": []string{strings.Repeat("a", consts.MaxTagNameLen+1)}},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(test_utils.Format(consts.MaxError, "Tags[0]", consts.MaxTagNameLen)),
	},
	{
		testBody:         JSON{"chall_id": "", "type": "aaa"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(test_utils.Format(consts.OneOfError, "Type", consts.DeployTypesStr)),
	},
	{
		testBody:         JSON{"chall_id": "", "max_points": -1},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(test_utils.Format(consts.MinError, "MaxPoints", 0)),
	},
	{
		testBody:         JSON{"chall_id": "", "score_type": "aaa"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(test_utils.Format(consts.OneOfError, "ScoreType", consts.ScoreTypesStr)),
	},
	{
		testBody:         JSON{"chall_id": "", "port": consts.MinPort - 1},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(test_utils.Format(consts.MinError, "Port", consts.MinPort)),
	},
	{
		testBody:         JSON{"chall_id": "", "port": consts.MaxPort + 1},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(test_utils.Format(consts.MaxError, "Port", consts.MaxPort)),
	},
	{
		testBody:         JSON{"chall_id": "", "conn_type": "aaa"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(test_utils.Format(consts.OneOfError, "ConnType", consts.ConnTypesStr)),
	},
	{
		testBody:         JSON{"chall_id": "", "lifetime": -1},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(test_utils.Format(consts.MinError, "Lifetime", 0)),
	},
	{
		testBody:         JSON{"chall_id": "", "lifetime": math.MaxInt32 + 1},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidJSON),
	},
	{
		testBody:         JSON{"chall_id": "", "envs": "<invalid-json>"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidEnvs),
	},
	{
		testBody:         JSON{"chall_id": "", "max_memory": -1},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(test_utils.Format(consts.MinError, "MaxMemory", 0)),
	},
	{
		testBody:         JSON{"chall_id": "", "max_memory": math.MaxInt32 + 1},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidJSON),
	},
	{
		testBody:         JSON{"chall_id": "", "max_cpu": "<invalid-float>"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidMaxCpu),
	},
	{
		testBody:         JSON{"chall_id": "", "max_cpu": "0.0"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidMaxCpu),
	},
	{
		testBody:         JSON{"chall_id": "", "max_cpu": "-1.0"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidMaxCpu),
	},
	{
		testBody:         JSON{"chall_id": "", "max_cpu": fmt.Sprintf("%d.0", math.MaxInt32+1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidMaxCpu),
	},
	{
		testBody:         JSON{"chall_id": -1, "name": "test"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(test_utils.Format(consts.MinError, "ChallID", 0)),
	},
	{
		testBody:         JSON{"chall_id": 9999, "name": "test"},
		expectedStatus:   http.StatusNotFound,
		expectedResponse: errorf(consts.ChallengeNotFound),
	},
	{
		testBody:         JSON{"chall_id": math.MaxInt32 + 1},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidJSON),
	},
	{
		testBody:         JSON{"chall_id": "", "name": "chall-2"},
		expectedStatus:   http.StatusConflict,
		expectedResponse: errorf(consts.ChallengeNameAlreadyExists),
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
			"conn_type":   "TCP_TLS",
			"description": "new test desc",
			"difficulty":  "Ez",
			"authors":     []string{"author1", "author2"},
			"tags":        []string{"tag1", "tag2", "tag3"},
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
		expectedStatus: http.StatusOK,
	},
}

func TestRoute(t *testing.T) {
	app := api.SetupApp(t.Context())
	defer api.Shutdown(app)

	test_utils.RegisterUser(t, "author", "author@test.test", "authorpass", sqlc.UserRoleAuthor)

	var challID int32
	session := test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "author@test.test", "password": "authorpass"}, http.StatusOK)
	session.Post("/categories", JSON{"name": "cat"}, -1)

	chall := test_utils.TryCreateChallenge(t, "chall", "cat", "test-desc", sqlc.DeployTypeNormal, 1, sqlc.ScoreTypeStatic)
	if chall != nil {
		challID = chall.ID
	}

	session.Patch("/challenges", nil, http.StatusBadRequest)
	session.CheckResponse(errorf(consts.InvalidJSON))

	for i, test := range testData {
		session := test_utils.NewApiTestSession(t, app)
		session.Post("/login", JSON{"email": "author@test.test", "password": "authorpass"}, http.StatusOK)

		if test.testBody != nil {
			if content, ok := test.testBody["chall_id"]; ok && content == "" {
				test.testBody["chall_id"] = challID
			}
		}

		session.Patch("/challenges", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)

		if i == len(testData)-1 {
			session.Get("/challenges", nil, http.StatusOK)
			body := session.Body()
			if body == nil {
				t.Fatal("Expected body to not be nil")
			}
			expected := JSON{
				"attachments": []string{},
				"authors":     test.testBody["authors"],
				"category":    test.testBody["category"],
				"conn_type":   test.testBody["conn_type"],
				"description": test.testBody["description"],
				"difficulty":  test.testBody["difficulty"],
				"first_blood": false,
				"hidden":      test.testBody["hidden"],
				"host":        test.testBody["host"],
				"id":          challID,
				"instance":    test.testBody["type"] != "Normal",
				"max_points":  test.testBody["max_points"],
				"name":        test.testBody["name"],
				"points":      test.testBody["max_points"],
				"port":        1234,
				"score_type":  test.testBody["score_type"],
				"solved":      false,
				"solves":      0,
				"tags":        test.testBody["tags"],
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

			session.Get(fmt.Sprintf("/challenges/%d", challID), nil, http.StatusOK)
			body = session.Body()
			if body == nil {
				t.Fatal("Expected body to not be nil")
			}
			expected = JSON{
				"docker_config": JSON{
					"compose":     test.testBody["compose"],
					"envs":        test.testBody["envs"],
					"hash_domain": test.testBody["hash_domain"],
					"image":       test.testBody["image"],
					"lifetime":    test.testBody["lifetime"],
					"max_cpu":     test.testBody["max_cpu"],
					"max_memory":  test.testBody["max_memory"],
				},
				"flags":       []string{},
				"solves_list": []string{},
				"type":        test.testBody["type"],
			}
			test_utils.Compare(t, expected, body)
		}
	}

	testBody := JSON{
		"chall_id":    challID,
		"name":        "Test",
		"category":    "cat-2",
		"description": "",
		"difficulty":  "",
		"authors":     []string{},
		"tags":        []string{},
		"hidden":      false,
		"max_points":  0,
		"host":        "",
		"port":        0,

		"image":       "",
		"compose":     "",
		"hash_domain": false,
		"lifetime":    0,
		"envs":        "",
		"max_memory":  0,
		"max_cpu":     "",
	}

	session = test_utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "author@test.test", "password": "authorpass"}, http.StatusOK)
	session.Patch("/challenges", testBody, http.StatusOK)
	session.CheckResponse(nil)

	session.Get("/challenges/", nil, http.StatusOK)
	body := session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	expected := JSON{
		"attachments": []string{},
		"authors":     testBody["authors"],
		"category":    testBody["category"],
		"conn_type":   "TCP_TLS",
		"description": testBody["description"],
		"difficulty":  testBody["difficulty"],
		"first_blood": false,
		"hidden":      testBody["hidden"],
		"host":        testBody["host"],
		"id":          challID,
		"instance":    testBody["type"] != "Normal",
		"max_points":  testBody["max_points"],
		"name":        testBody["name"],
		"points":      testBody["max_points"],
		"port":        testBody["port"],
		"score_type":  "Dynamic",
		"solved":      false,
		"solves":      0,
		"tags":        testBody["tags"],
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

	session.Get(fmt.Sprintf("/challenges/%d", challID), nil, http.StatusOK)
	body = session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	expected = JSON{
		"docker_config": JSON{
			"compose":     testBody["compose"],
			"envs":        testBody["envs"],
			"hash_domain": testBody["hash_domain"],
			"image":       testBody["image"],
			"lifetime":    testBody["lifetime"],
			"max_cpu":     testBody["max_cpu"],
			"max_memory":  testBody["max_memory"],
		},
		"flags":       []string{},
		"solves_list": []string{},
		"type":        "Container",
	}
	test_utils.Compare(t, expected, body)
}
