package test_utils

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"trxd/api/routes/challenges_create"
	"trxd/api/routes/teams_register"
	"trxd/api/routes/users_register"
	"trxd/db"
	"trxd/db/sqlc"
	"trxd/utils"
)

const PROJECT_DIR = "backend"

func fatalf(format string, a ...any) {
	fmt.Printf(format, a...)
	os.Exit(1)
}

func Main(m *testing.M) {
	dir, err := os.Getwd()
	if err != nil {
		fatalf("Failed to get current directory: %v\n", err)
	}

	parts := strings.Split(dir, string(os.PathSeparator))
	depth := -1
	for i, part := range parts {
		if part == PROJECT_DIR {
			depth = len(parts) - (i + 1)
			break
		}
	}

	if depth < 0 {
		fatalf("Failed to determine depth from directory: %s\n", dir)
	}

	err = os.Chdir(strings.Repeat("../", depth))
	if err != nil {
		fatalf("Failed to change directory: %v\n", err)
	}

	name := filepath.Base(dir)
	err = db.OpenTestDB("test_" + name)
	if err != nil {
		fatalf("Failed to open test database: %v\n", err)
	}
	defer db.CloseTestDB()

	err = db.DeleteAll()
	if err != nil {
		fatalf("Failed to delete all data: %v\n", err)
	}

	err = db.InitConfigs()
	if err != nil {
		fatalf("Failed to initialize configs: %v\n", err)
	}

	err = db.InsertMockData()
	if err != nil {
		fatalf("Failed to insert mock data: %v\n", err)
	}

	err = db.UpdateConfig(context.Background(), "allow-register", "true")
	if err != nil {
		fatalf("Failed to update config: %v\n", err)
	}

	exitCode := m.Run()
	os.Exit(exitCode)
}

func UpdateConfig(t *testing.T, name string, value string) {
	err := db.UpdateConfig(t.Context(), name, value)
	if err != nil {
		t.Fatalf("Failed to update config %s: %v", name, err)
	}
}

func RegisterUser(t *testing.T, name, email, password string, role sqlc.UserRole) *sqlc.User {
	user, err := users_register.RegisterUser(t.Context(), name, email, password, role)
	if err != nil {
		t.Fatalf("Failed to register user %s: %v", name, err)
	}
	if user == nil {
		t.Fatalf("Registered user '%s' is nil", name)
	}

	return user
}

func RegisterTeam(t *testing.T, name, password string, userID int32) *sqlc.Team {
	team, err := teams_register.RegisterTeam(t.Context(), name, password, userID)
	if err != nil {
		t.Fatalf("Failed to register team %s: %v", name, err)
	}
	if team == nil {
		t.Fatalf("Registered team '%s' is nil", name)
	}

	return team
}

func CreateChallenge(t *testing.T, name string, category string, description string,
	challType sqlc.DeployType, maxPoints int32, scoreType sqlc.ScoreType) *sqlc.Challenge {
	chall, err := challenges_create.CreateChallenge(t.Context(), name, category, description, challType, maxPoints, scoreType)
	if err != nil {
		t.Fatalf("Failed to create challenge %s: %v", name, err)
	}
	if chall == nil {
		t.Fatalf("Challenge creation of '%s' returned nil", name)
	}

	return chall
}

func TryCreateChallenge(t *testing.T, name string, category string, description string,
	challType sqlc.DeployType, maxPoints int32, scoreType sqlc.ScoreType) *sqlc.Challenge {
	chall, err := challenges_create.CreateChallenge(t.Context(), name, category, description, challType, maxPoints, scoreType)
	if err != nil {
		t.Fatalf("Failed to create challenge %s: %v", name, err)
	}

	return chall
}

func GetTeamByName(t *testing.T, name string) *sqlc.Team {
	team, err := db.GetTeamByName(t.Context(), name)
	if err != nil {
		t.Fatalf("Failed to get team %s: %v", name, err)
	}
	if team == nil {
		t.Fatalf("Team %s not found", name)
	}

	return team
}

func Compare(t *testing.T, a, b interface{}) {
	err := utils.Compare(a, b)
	if err != nil {
		t.Fatalf("Failed to compare values: %v", err)
	}
}

func DeleteKeys(data interface{}, keys ...string) interface{} {
	switch val := data.(type) {
	case map[string]interface{}:
		for _, key := range keys {
			delete(val, key)
		}
		for k, v := range val {
			val[k] = DeleteKeys(v, keys...)
		}
		return val

	case []interface{}:
		for i, v := range val {
			val[i] = DeleteKeys(v, keys...)
		}
		return val

	default:
		return val
	}
}

func GetModuleName(t *testing.T) string {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	return filepath.Base(dir)
}

func CreateDir(t *testing.T, dir string) {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		t.Fatalf("Failed to create directory %s: %v", dir, err)
	}
}

func CreateFile(t *testing.T, file string, content string) {
	f, err := os.Create(file)
	if err != nil {
		t.Fatalf("Failed to create file %s: %v", file, err)
	}
	defer f.Close()

	_, err = f.WriteString(content)
	if err != nil {
		t.Fatalf("Failed to write content to file %s: %v", file, err)
	}
}
