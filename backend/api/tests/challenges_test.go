package tests

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"trxd/api"
	"trxd/db"
	"trxd/utils"
	"trxd/utils/consts"
)

var testCreateChallenge = []struct {
	testBody         interface{}
	expectedStatus   int
	expectedResponse JSON
}{
	{
		testBody:         nil,
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidJSON),
	},
	{
		testBody:         JSON{"name": "test"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"category": "test"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"description": "test"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"type": "test"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"max_points": 1},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"score_type": "test"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"name": strings.Repeat("a", consts.MaxChallNameLength+1), "category": "cat", "description": "test-desc", "type": "Normal", "max_points": 1, "score_type": "Static"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.LongChallName),
	},
	{
		testBody:         JSON{"name": "test", "category": strings.Repeat("a", consts.MaxCategoryLength+1), "description": "test-desc", "type": "Normal", "max_points": 1, "score_type": "Static"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.LongCategory),
	},
	{
		testBody:         JSON{"name": "test", "category": "cat", "description": strings.Repeat("a", consts.MaxChallDescLength+1), "type": "Normal", "max_points": 1, "score_type": "Static"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.LongChallDesc),
	},
	{
		testBody:         JSON{"name": "test", "category": "cat", "description": "test-desc", "type": "aaaaa", "max_points": 1, "score_type": "Static"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidChallType),
	},
	{
		testBody:         JSON{"name": "test", "category": "cat", "description": "test-desc", "type": "Normal", "max_points": 0, "score_type": "Static"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidChallMaxPoints),
	},
	{
		testBody:         JSON{"name": "test", "category": "cat", "description": "test-desc", "type": "Normal", "max_points": 1, "score_type": "aaaa"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidChallScoreType),
	},
	{
		testBody:         JSON{"name": "test3", "category": "cat2", "description": "test-desc", "type": "Normal", "max_points": 1, "score_type": "Static"},
		expectedStatus:   http.StatusNotFound,
		expectedResponse: errorf(consts.CategoryNotFound),
	},
	{
		testBody:         JSON{"name": "test", "category": "cat", "description": "test-desc", "type": "Normal", "max_points": 1, "score_type": "Static"},
		expectedStatus:   http.StatusOK,
		expectedResponse: JSON{"name": "test", "category": "cat"},
	},
	{
		testBody:         JSON{"name": "test", "category": "cat", "description": "test-desc", "type": "Normal", "max_points": 1, "score_type": "Static"},
		expectedStatus:   http.StatusConflict,
		expectedResponse: errorf(consts.ChallengeAlreadyExists),
	},
	{
		testBody:         JSON{"name": "test2", "category": "cat", "description": "test-desc", "type": "Normal", "max_points": 1, "score_type": "Static"},
		expectedStatus:   http.StatusOK,
		expectedResponse: JSON{"name": "test2", "category": "cat"},
	},
}

func TestCreateChallenge(t *testing.T) {
	db.DeleteAll()
	app := api.SetupApp()
	defer app.Shutdown()

	_, err := db.RegisterUser(context.Background(), "author", "author@test.test", "authorpass", db.UserRoleAuthor)
	if err != nil {
		t.Fatalf("Failed to register author user: %v", err)
	}
	_, err = db.CreateCategory(context.Background(), "cat", "icon")
	if err != nil {
		t.Fatalf("Failed to create category: %v", err)
	}

	for _, test := range testCreateChallenge {
		session := utils.NewApiTestSession(t, app)
		session.Post("/login", JSON{"email": "author@test.test", "password": "authorpass"}, http.StatusOK)
		session.Post("/challenges", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}
}

var testCreateFlag = []struct {
	testBody         interface{}
	expectedStatus   int
	expectedResponse JSON
}{
	{
		testBody:         nil,
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidJSON),
	},
	{
		testBody:         JSON{"chall_id": ""},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"flag": "test"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"chall_id": "", "flag": strings.Repeat("a", consts.MaxFlagLength+1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.LongFlag),
	},
	{
		testBody:         JSON{"chall_id": 99999, "flag": "flag{test}"},
		expectedStatus:   http.StatusNotFound,
		expectedResponse: errorf(consts.ChallengeNotFound),
	},
	{
		testBody:       JSON{"chall_id": "", "flag": "test"},
		expectedStatus: http.StatusOK,
	},
	{
		testBody:       JSON{"chall_id": "", "flag": `flag\{test\}`, "regex": true},
		expectedStatus: http.StatusOK,
	},
	{
		testBody:         JSON{"chall_id": "", "flag": "test"},
		expectedStatus:   http.StatusConflict,
		expectedResponse: errorf(consts.FlagAlreadyExists),
	},
}

func TestCreateFlag(t *testing.T) {
	db.DeleteAll()
	app := api.SetupApp()
	defer app.Shutdown()

	user, err := db.RegisterUser(context.Background(), "test", "test@test.test", "testpass", db.UserRoleAuthor)
	if err != nil {
		t.Fatalf("Failed to register author user: %v", err)
	}
	if user == nil {
		t.Fatal("User registration returned nil")
	}

	cat, err := db.CreateCategory(context.Background(), "cat", "icon")
	if err != nil {
		t.Fatalf("Failed to create category: %v", err)
	}
	if cat == nil {
		t.Fatal("Category creation returned nil")
	}
	chall, err := db.CreateChallenge(context.Background(), "chall", cat.Name, "test-desc", db.DeployTypeNormal, 1, db.ScoreTypeStatic)
	if err != nil {
		t.Fatalf("Failed to create challenge: %v", err)
	}
	if chall == nil {
		t.Fatal("Challenge creation returned nil")
	}

	for _, test := range testCreateFlag {
		session := utils.NewApiTestSession(t, app)
		session.Post("/login", JSON{"email": "test@test.test", "password": "testpass"}, http.StatusOK)
		if body, ok := test.testBody.(JSON); ok && body != nil {
			if content, ok := body["chall_id"]; ok && content == "" {
				test.testBody.(JSON)["chall_id"] = chall.ID
			}
		}
		session.Post("/flag", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}
}

var testSubmit = []struct {
	testBody         interface{}
	expectedStatus   int
	expectedResponse JSON
}{
	{
		testBody:         nil,
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidJSON),
	},
	{
		testBody:         JSON{"chall_id": ""},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"flag": "test"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"chall_id": "", "flag": strings.Repeat("a", consts.MaxFlagLength+1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.LongFlag),
	},
	{
		testBody:         JSON{"chall_id": 99999, "flag": "flag{test}"},
		expectedStatus:   http.StatusNotFound,
		expectedResponse: errorf(consts.ChallengeNotFound),
	},
	{
		testBody:         JSON{"chall_id": "", "flag": "test"},
		expectedStatus:   http.StatusOK,
		expectedResponse: JSON{"status": string(db.SubmissionStatusWrong)},
	},
	{
		testBody:         JSON{"chall_id": "", "flag": "flag{test}"},
		expectedStatus:   http.StatusOK,
		expectedResponse: JSON{"status": string(db.SubmissionStatusCorrect)},
	},
	{
		testBody:         JSON{"chall_id": "", "flag": "flag{test}"},
		expectedStatus:   http.StatusOK,
		expectedResponse: JSON{"status": string(db.SubmissionStatusRepeated)},
	},
}

func TestSubmit(t *testing.T) {
	db.DeleteAll()
	db.InitConfigs()
	app := api.SetupApp()
	defer app.Shutdown()

	err := db.UpdateConfig(context.Background(), "allow-register", "true")
	if err != nil {
		t.Fatalf("Failed to update config: %v", err)
	}

	session := utils.NewApiTestSession(t, app)
	session.Post("/register", JSON{"username": "test2", "email": "test2@test.test", "password": "testpass"}, http.StatusOK)
	session.Post("/submit", JSON{"chall_id": 0, "flag": "flag{test}"}, http.StatusForbidden)
	session.CheckResponse(errorf(consts.Forbidden))

	user3, err := db.RegisterUser(context.Background(), "test3", "test3@test.test", "testpass", db.UserRoleAdmin)
	if err != nil {
		t.Fatalf("Failed to register test user: %v", err)
	}
	if user3 == nil {
		t.Fatal("User registration returned nil")
	}
	session = utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "test3@test.test", "password": "testpass"}, http.StatusOK)
	session.Post("/submit", JSON{"chall_id": 0, "flag": "flag{test}"}, http.StatusNotFound)
	session.CheckResponse(errorf(consts.ChallengeNotFound))

	user, err := db.RegisterUser(context.Background(), "test", "test@test.test", "testpass")
	if err != nil {
		t.Fatalf("Failed to register test user: %v", err)
	}
	if user == nil {
		t.Fatal("User registration returned nil")
	}
	team, err := db.RegisterTeam(context.Background(), "test-team", "teampasswd", user.ID)
	if err != nil {
		t.Fatalf("Failed to register test team: %v", err)
	}
	if team == nil {
		t.Fatal("Team registration returned nil")
	}

	cat, err := db.CreateCategory(context.Background(), "cat", "icon")
	if err != nil {
		t.Fatalf("Failed to create category: %v", err)
	}
	if cat == nil {
		t.Fatal("Category creation returned nil")
	}
	chall, err := db.CreateChallenge(context.Background(), "chall", cat.Name, "test-desc", db.DeployTypeNormal, 1, db.ScoreTypeDynamic)
	if err != nil {
		t.Fatalf("Failed to create challenge: %v", err)
	}
	if chall == nil {
		t.Fatal("Challenge creation returned nil")
	}
	flag, err := db.CreateFlag(context.Background(), chall.ID, "flag{test}", false)
	if err != nil {
		t.Fatalf("Failed to create flag: %v", err)
	}
	if flag == nil {
		t.Fatal("Flag creation returned nil")
	}

	for _, test := range testSubmit {
		session := utils.NewApiTestSession(t, app)
		session.Post("/login", JSON{"email": "test@test.test", "password": "testpass"}, http.StatusOK)
		if body, ok := test.testBody.(JSON); ok && body != nil {
			if content, ok := body["chall_id"]; ok && content == "" {
				test.testBody.(JSON)["chall_id"] = chall.ID
			}
		}
		session.Post("/submit", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}
}

var testDeleteChallenge = []struct {
	testBody         interface{}
	expectedStatus   int
	expectedResponse JSON
}{
	{
		testBody:         nil,
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidJSON),
	},
	{
		testBody:         JSON{},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"chall_id": -1},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidChallengeID),
	},
	{
		testBody:       JSON{"chall_id": ""},
		expectedStatus: http.StatusOK,
	},
	{
		testBody:       JSON{"chall_id": ""},
		expectedStatus: http.StatusOK,
	},
}

func TestDeleteChallenge(t *testing.T) {
	db.DeleteAll()
	app := api.SetupApp()
	defer app.Shutdown()

	_, err := db.RegisterUser(context.Background(), "author", "author@test.test", "authorpass", db.UserRoleAuthor)
	if err != nil {
		t.Fatalf("Failed to register author user: %v", err)
	}

	var challID int32
	for _, test := range testDeleteChallenge {
		_, err := db.CreateCategory(context.Background(), "cat", "icon")
		if err != nil {
			t.Fatalf("Failed to create category: %v", err)
		}
		chall, err := db.CreateChallenge(context.Background(), "chall", "cat", "test-desc", db.DeployTypeNormal, 1, db.ScoreTypeStatic)
		if err != nil {
			t.Fatalf("Failed to create challenge: %v", err)
		}
		if chall != nil {
			challID = chall.ID
		}

		session := utils.NewApiTestSession(t, app)
		session.Post("/login", JSON{"email": "author@test.test", "password": "authorpass"}, http.StatusOK)
		if body, ok := test.testBody.(JSON); ok && body != nil {
			if content, ok := body["chall_id"]; ok && content == "" {
				test.testBody.(JSON)["chall_id"] = challID
			}
		}
		session.Delete("/challenges", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}
}

var testDeleteFlag = []struct {
	testBody         interface{}
	expectedStatus   int
	expectedResponse JSON
}{
	{
		testBody:         nil,
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.InvalidJSON),
	},
	{
		testBody:         JSON{"chall_id": ""},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"flag": "test"},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.MissingRequiredFields),
	},
	{
		testBody:         JSON{"chall_id": "", "flag": strings.Repeat("a", consts.MaxFlagLength+1)},
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: errorf(consts.LongFlag),
	},
	{
		testBody:         JSON{"chall_id": 99999, "flag": "flag{test}"},
		expectedStatus:   http.StatusNotFound,
		expectedResponse: errorf(consts.ChallengeNotFound),
	},
	{
		testBody:       JSON{"chall_id": "", "flag": "test"},
		expectedStatus: http.StatusOK,
	},
	{
		testBody:       JSON{"chall_id": "", "flag": "test"},
		expectedStatus: http.StatusOK,
	},
}

func TestDeleteFlag(t *testing.T) {
	db.DeleteAll()
	app := api.SetupApp()
	defer app.Shutdown()

	user, err := db.RegisterUser(context.Background(), "test", "test@test.test", "testpass", db.UserRoleAuthor)
	if err != nil {
		t.Fatalf("Failed to register author user: %v", err)
	}
	if user == nil {
		t.Fatal("User registration returned nil")
	}

	cat, err := db.CreateCategory(context.Background(), "cat", "icon")
	if err != nil {
		t.Fatalf("Failed to create category: %v", err)
	}
	if cat == nil {
		t.Fatal("Category creation returned nil")
	}
	chall, err := db.CreateChallenge(context.Background(), "chall", cat.Name, "test-desc", db.DeployTypeNormal, 1, db.ScoreTypeStatic)
	if err != nil {
		t.Fatalf("Failed to create challenge: %v", err)
	}
	if chall == nil {
		t.Fatal("Challenge creation returned nil")
	}

	for _, test := range testDeleteFlag {
		_, err := db.CreateFlag(context.Background(), chall.ID, "test", false)
		if err != nil {
			t.Fatalf("Failed to create flag: %v", err)
		}

		session := utils.NewApiTestSession(t, app)
		session.Post("/login", JSON{"email": "test@test.test", "password": "testpass"}, http.StatusOK)
		if body, ok := test.testBody.(JSON); ok && body != nil {
			if content, ok := body["chall_id"]; ok && content == "" {
				test.testBody.(JSON)["chall_id"] = chall.ID
			}
		}
		session.Delete("/flag", test.testBody, test.expectedStatus)
		session.CheckResponse(test.expectedResponse)
	}
}

func TestGetChallenges(t *testing.T) {
	db.DeleteAll()
	db.InitConfigs()
	db.InsertMockData()
	app := api.SetupApp()
	defer app.Shutdown()

	err := db.UpdateConfig(context.Background(), "allow-register", "true")
	if err != nil {
		t.Fatalf("Failed to update config: %v", err)
	}

	expectedPlayer := []JSON{
		{
			"attachments": []interface{}{},
			"authors": []interface{}{
				"author1",
				"author2",
			},
			"category":    "cat-1",
			"description": "TEST chall-1 DESC",
			"difficulty":  "Easy",
			"flags":       nil,
			"hidden":      false,
			"host":        "http://theromanxpl0.it",
			"instance":    false,
			"name":        "chall-1",
			"points":      500,
			"port":        1337,
			"solved":      false,
			"solves":      1,
			"solves_list": []map[string]interface{}{
				{
					"name": "A",
				},
			},
			"tags": []interface{}{
				"tag-1",
				"test-tag",
			},
			"timeout": 0,
		},
		{
			"attachments": []interface{}{},
			"authors": []interface{}{
				"author1",
			},
			"category":    "cat-1",
			"description": "TEST chall-3 DESC",
			"difficulty":  "Hard",
			"flags":       nil,
			"hidden":      false,
			"host":        "",
			"instance":    true,
			"name":        "chall-3",
			"points":      500,
			"port":        0,
			"solved":      false,
			"solves":      1,
			"solves_list": []map[string]interface{}{
				{
					"name": "A",
				},
			},
			"tags": []interface{}{
				"tag-3",
			},
			"timeout": 0,
		},
		{
			"attachments": []interface{}{},
			"authors": []interface{}{
				"author2",
			},
			"category":    "cat-1",
			"description": "TEST chall-4 DESC",
			"difficulty":  "Insane",
			"flags":       nil,
			"hidden":      false,
			"host":        "",
			"instance":    true,
			"name":        "chall-4",
			"points":      498,
			"port":        0,
			"solved":      false,
			"solves":      2,
			"solves_list": []map[string]interface{}{
				{
					"name": "A",
				},
				{
					"name": "B",
				},
			},
			"tags": []interface{}{
				"tag-4",
			},
			"timeout": 0,
		},
		{
			"attachments": []interface{}{},
			"authors": []interface{}{
				"author1",
				"author2",
				"author3",
			},
			"category":    "cat-2",
			"description": "TEST chall-2 DESC",
			"difficulty":  "Medium",
			"flags":       nil,
			"hidden":      false,
			"host":        "",
			"instance":    false,
			"name":        "chall-2",
			"points":      500,
			"port":        0,
			"solved":      false,
			"solves":      1,
			"solves_list": []map[string]interface{}{
				{
					"name": "B",
				},
			},
			"tags": []interface{}{
				"tag-2",
			},
			"timeout": 0,
		},
	}

	session := utils.NewApiTestSession(t, app)
	session.Post("/register", JSON{"username": "test", "email": "test2@test.test", "password": "testpass"}, http.StatusOK)
	session.Post("/teams", JSON{"name": "test-team", "password": "testpass"}, http.StatusOK)
	session.Get("/challenges", nil, http.StatusOK)
	body := session.Body()
	for _, chall := range body.([]interface{}) {
		delete(chall.(map[string]interface{}), "id")
		for _, solve := range chall.(map[string]interface{})["solves_list"].([]interface{}) {
			delete(solve.(map[string]interface{}), "id")
			delete(solve.(map[string]interface{}), "timestamp")
		}
	}
	err = utils.Compare(expectedPlayer, body)
	if err != nil {
		t.Fatalf("Compare Error: %v", err)
	}

	expectedAuthor := []JSON{
		{
			"attachments": []interface{}{},
			"authors": []interface{}{
				"author3",
			},
			"category":    "cat-2",
			"description": "TEST chall-5 DESC",
			"difficulty":  "Easy",
			"flags": []map[string]interface{}{
				{
					"flag":  "flag{test-5}",
					"regex": false,
				},
			},
			"hidden":      true,
			"host":        "",
			"instance":    false,
			"name":        "chall-5",
			"points":      500,
			"port":        0,
			"solved":      false,
			"solves":      0,
			"solves_list": []map[string]interface{}{},
			"tags": []interface{}{
				"tag-5",
			},
			"timeout": 0,
		},
		{
			"attachments": []interface{}{},
			"authors": []interface{}{
				"author1",
				"author2",
			},
			"category":    "cat-1",
			"description": "TEST chall-1 DESC",
			"difficulty":  "Easy",
			"flags": []map[string]interface{}{
				{
					"flag":  "flag{test-1}",
					"regex": false,
				},
				{
					"flag":  "flag\\{test-[a-z]{2}\\}",
					"regex": true,
				},
			},
			"hidden":   false,
			"host":     "http://theromanxpl0.it",
			"instance": false,
			"name":     "chall-1",
			"points":   500,
			"port":     1337,
			"solved":   false,
			"solves":   1,
			"solves_list": []map[string]interface{}{
				{
					"name": "A",
				},
			},
			"tags": []interface{}{
				"tag-1",
				"test-tag",
			},
			"timeout": 0,
		},
		{
			"attachments": []interface{}{},
			"authors": []interface{}{
				"author1",
			},
			"category":    "cat-1",
			"description": "TEST chall-3 DESC",
			"difficulty":  "Hard",
			"flags": []map[string]interface{}{
				{
					"flag":  "flag{test-3}",
					"regex": false,
				},
			},
			"hidden":   false,
			"host":     "",
			"instance": true,
			"name":     "chall-3",
			"points":   500,
			"port":     0,
			"solved":   false,
			"solves":   1,
			"solves_list": []map[string]interface{}{
				{
					"name": "A",
				},
			},
			"tags": []interface{}{
				"tag-3",
			},
			"timeout": 0,
		},
		{
			"attachments": []interface{}{},
			"authors": []interface{}{
				"author2",
			},
			"category":    "cat-1",
			"description": "TEST chall-4 DESC",
			"difficulty":  "Insane",
			"flags": []map[string]interface{}{
				{
					"flag":  "flag{test-4}",
					"regex": false,
				},
			},
			"hidden":   false,
			"host":     "",
			"instance": true,
			"name":     "chall-4",
			"points":   498,
			"port":     0,
			"solved":   false,
			"solves":   2,
			"solves_list": []map[string]interface{}{
				{
					"name": "A",
				},
				{
					"name": "B",
				},
			},
			"tags": []interface{}{
				"tag-4",
			},
			"timeout": 0,
		},
		{
			"attachments": []interface{}{},
			"authors": []interface{}{
				"author1",
				"author2",
				"author3",
			},
			"category":    "cat-2",
			"description": "TEST chall-2 DESC",
			"difficulty":  "Medium",
			"flags": []map[string]interface{}{
				{
					"flag":  "flag{test-2}",
					"regex": false,
				},
			},
			"hidden":   false,
			"host":     "",
			"instance": false,
			"name":     "chall-2",
			"points":   500,
			"port":     0,
			"solved":   false,
			"solves":   1,
			"solves_list": []map[string]interface{}{
				{
					"name": "B",
				},
			},
			"tags": []interface{}{
				"tag-2",
			},
			"timeout": 0,
		},
	}

	user, err := db.RegisterUser(context.Background(), "test2", "test3@test.test", "testpass", db.UserRoleAuthor)
	if err != nil {
		t.Fatalf("Failed to register author user: %v", err)
	}
	if user == nil {
		t.Fatal("User registration returned nil")
	}

	session = utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "test3@test.test", "password": "testpass"}, http.StatusOK)
	session.Get("/challenges", nil, http.StatusOK)
	body = session.Body()
	for _, chall := range body.([]interface{}) {
		delete(chall.(map[string]interface{}), "id")
		for _, solve := range chall.(map[string]interface{})["solves_list"].([]interface{}) {
			delete(solve.(map[string]interface{}), "id")
			delete(solve.(map[string]interface{}), "timestamp")
		}
	}
	err = utils.Compare(expectedAuthor, body)
	if err != nil {
		t.Fatalf("Compare Error: %v", err)
	}
}

func TestGetChallenge(t *testing.T) {
	db.DeleteAll()
	db.InitConfigs()
	db.InsertMockData()
	app := api.SetupApp()
	defer app.Shutdown()

	err := db.UpdateConfig(context.Background(), "allow-register", "true")
	if err != nil {
		t.Fatalf("Failed to update config: %v", err)
	}

	session := utils.NewApiTestSession(t, app)
	session.Post("/register", JSON{"username": "test", "email": "test2@test.test", "password": "testpass"}, http.StatusOK)
	session.Post("/teams", JSON{"name": "test-team", "password": "testpass"}, http.StatusOK)
	session.Get("/challenges", nil, http.StatusOK)
	body := session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	id := int32(body.([]interface{})[0].(map[string]interface{})["id"].(float64))

	expectedPlayer := JSON{
		"attachments": []interface{}{},
		"authors": []interface{}{
			"author1",
			"author2",
		},
		"category":    "cat-1",
		"description": "TEST chall-1 DESC",
		"difficulty":  "Easy",
		"flags":       nil,
		"hidden":      false,
		"host":        "http://theromanxpl0.it",
		"instance":    false,
		"name":        "chall-1",
		"points":      500,
		"port":        1337,
		"solved":      false,
		"solves":      1,
		"solves_list": []map[string]interface{}{
			{
				"name": "A",
			},
		},
		"tags": []interface{}{
			"tag-1",
			"test-tag",
		},
		"timeout": 0,
	}

	session = utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "test2@test.test", "password": "testpass"}, http.StatusOK)
	session.Get(fmt.Sprintf("/challenges/%d", id), nil, http.StatusOK)
	body = session.Body()
	delete(body.(map[string]interface{}), "id")
	for _, solve := range body.(map[string]interface{})["solves_list"].([]interface{}) {
		delete(solve.(map[string]interface{}), "id")
		delete(solve.(map[string]interface{}), "timestamp")
	}
	err = utils.Compare(expectedPlayer, body)
	if err != nil {
		t.Fatalf("Compare Error: %v", err)
	}

	expectedAuthor := JSON{
		"attachments": []interface{}{},
		"authors": []interface{}{
			"author1",
			"author2",
		},
		"category":    "cat-1",
		"description": "TEST chall-1 DESC",
		"difficulty":  "Easy",
		"flags": []map[string]interface{}{
			{
				"flag":  "flag{test-1}",
				"regex": false,
			},
			{
				"flag":  "flag\\{test-[a-z]{2}\\}",
				"regex": true,
			},
		},
		"hidden":   false,
		"host":     "http://theromanxpl0.it",
		"instance": false,
		"name":     "chall-1",
		"points":   500,
		"port":     1337,
		"solved":   false,
		"solves":   1,
		"solves_list": []map[string]interface{}{
			{
				"name": "A",
			},
		},
		"tags": []interface{}{
			"tag-1",
			"test-tag",
		},
		"timeout": 0,
	}

	user, err := db.RegisterUser(context.Background(), "test2", "test3@test.test", "testpass", db.UserRoleAuthor)
	if err != nil {
		t.Fatalf("Failed to register author user: %v", err)
	}
	if user == nil {
		t.Fatal("User registration returned nil")
	}

	session = utils.NewApiTestSession(t, app)
	session.Post("/login", JSON{"email": "test3@test.test", "password": "testpass"}, http.StatusOK)
	session.Get(fmt.Sprintf("/challenges/%d", id), nil, http.StatusOK)
	body = session.Body()
	delete(body.(map[string]interface{}), "id")
	for _, solve := range body.(map[string]interface{})["solves_list"].([]interface{}) {
		delete(solve.(map[string]interface{}), "id")
		delete(solve.(map[string]interface{}), "timestamp")
	}
	err = utils.Compare(expectedAuthor, body)
	if err != nil {
		t.Fatalf("Compare Error: %v", err)
	}

	expectedAuthorHidden := JSON{
		"attachments": []interface{}{},
		"authors": []interface{}{
			"author3",
		},
		"category":    "cat-2",
		"description": "TEST chall-5 DESC",
		"difficulty":  "Easy",
		"flags": []map[string]interface{}{
			{
				"flag":  "flag{test-5}",
				"regex": false,
			},
		},
		"hidden":      true,
		"host":        "",
		"instance":    false,
		"name":        "chall-5",
		"points":      500,
		"port":        0,
		"solved":      false,
		"solves":      0,
		"solves_list": []map[string]interface{}{},
		"tags": []interface{}{
			"tag-5",
		},
		"timeout": 0,
	}

	session.Get("/challenges", nil, http.StatusOK)
	body = session.Body()
	if body == nil {
		t.Fatal("Expected body to not be nil")
	}
	id = int32(body.([]interface{})[0].(map[string]interface{})["id"].(float64))

	session.Get(fmt.Sprintf("/challenges/%d", id), nil, http.StatusOK)
	body = session.Body()
	delete(body.(map[string]interface{}), "id")
	for _, solve := range body.(map[string]interface{})["solves_list"].([]interface{}) {
		delete(solve.(map[string]interface{}), "id")
		delete(solve.(map[string]interface{}), "timestamp")
	}
	err = utils.Compare(expectedAuthorHidden, body)
	if err != nil {
		t.Fatalf("Compare Error: %v", err)
	}
}
