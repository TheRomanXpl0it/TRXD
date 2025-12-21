package test_utils

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"trxd/api/routes/challenges_create"
	"trxd/api/routes/challenges_update"
	"trxd/api/routes/teams_register"
	"trxd/api/routes/users_register"
	"trxd/db"
	"trxd/db/sqlc"
	"trxd/utils"
	"trxd/utils/consts"
	"trxd/utils/crypto_utils"
)

const PROJECT_DIR = "backend"

func fatalf(format string, a ...any) {
	fmt.Printf(format, a...)
	os.Exit(1)
}

func Main(m *testing.M) {
	ctx := context.Background()

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
	defer func() {
		err := db.CloseTestDB()
		if err != nil {
			fatalf("Failed to close test database: %v\n", err)
		}
	}()

	err = db.DeleteAll(ctx)
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

	consts.DefaultConfigs["allow-register"] = true
	err = db.UpdateConfig(ctx, "allow-register", "true")
	if err != nil {
		fatalf("Failed to update config: %v\n", err)
	}
	consts.DefaultConfigs["secret"] = "test-secret"
	err = db.UpdateConfig(ctx, "secret", "test-secret")
	if err != nil {
		fatalf("Failed to update config: %v\n", err)
	}
	consts.DefaultConfigs["domain"] = "test.com"
	err = db.UpdateConfig(ctx, "domain", "test.com")
	if err != nil {
		fatalf("Failed to update config: %v\n", err)
	}

	exitCode := m.Run()
	os.Exit(exitCode)
}

func Fatalf(t *testing.T, msg string, a ...any) {
	_, file, line, _ := runtime.Caller(2)
	t.Fatalf("%s:%d: %s", file, line, fmt.Sprintf(msg, a...))
}

func UpdateConfig(t *testing.T, name string, value string) {
	err := db.UpdateConfig(t.Context(), name, value)
	if err != nil {
		Fatalf(t, "Failed to update config %s: %v", name, err)
	}
}

func RegisterUser(t *testing.T, name, email, password string, role sqlc.UserRole) *sqlc.User {
	tx, err := db.BeginTx(t.Context())
	if err != nil {
		Fatalf(t, "Failed to begin transaction: %v", err)
	}
	defer db.Rollback(tx)

	user, err := users_register.RegisterUser(t.Context(), tx, name, email, password, role)
	if err != nil {
		Fatalf(t, "Failed to register user %s: %v", name, err)
	}
	if user == nil {
		Fatalf(t, "Registered user '%s' is nil", name)
	}

	err = tx.Commit()
	if err != nil {
		Fatalf(t, "Failed to commit transaction: %v", err)
	}

	return user
}

func RegisterTeam(t *testing.T, name, password string, userID int32) *sqlc.Team {
	tx, err := db.BeginTx(t.Context())
	if err != nil {
		Fatalf(t, "Failed to begin transaction: %v", err)
	}
	defer db.Rollback(tx)

	team, err := teams_register.RegisterTeam(t.Context(), tx, name, password, userID)
	if err != nil {
		Fatalf(t, "Failed to register team %s: %v", name, err)
	}
	if team == nil {
		Fatalf(t, "Registered team '%s' is nil", name)
	}

	err = tx.Commit()
	if err != nil {
		Fatalf(t, "Failed to commit transaction: %v", err)
	}

	return team
}

func CreateChallenge(t *testing.T, name string, category string, description string,
	challType sqlc.DeployType, maxPoints int32, scoreType sqlc.ScoreType) *sqlc.Challenge {
	chall, err := challenges_create.CreateChallenge(t.Context(), name, category, description, challType, maxPoints, scoreType)
	if err != nil {
		Fatalf(t, "Failed to create challenge %s: %v", name, err)
	}
	if chall == nil {
		Fatalf(t, "Challenge creation of '%s' returned nil", name)
	}

	return chall
}

func TryCreateChallenge(t *testing.T, name string, category string, description string,
	challType sqlc.DeployType, maxPoints int32, scoreType sqlc.ScoreType) *sqlc.Challenge {
	chall, err := challenges_create.CreateChallenge(t.Context(), name, category, description, challType, maxPoints, scoreType)
	if err != nil {
		Fatalf(t, "Failed to create challenge %s: %v", name, err)
	}

	return chall
}

func UnveilChallenge(t *testing.T, id int32) {
	False := false
	err := challenges_update.UpdateChallenge(t.Context(), &challenges_update.UpdateChallParams{
		ChallID: &id,
		Hidden:  &False,
	})
	if err != nil {
		Fatalf(t, "Failed to update challenge %d: %v", id, err)
	}
}

func GetTeamByName(t *testing.T, name string) *sqlc.Team {
	team, err := db.GetTeamByName(t.Context(), name)
	if err != nil {
		Fatalf(t, "Failed to get team %s: %v", name, err)
	}
	if team == nil {
		Fatalf(t, "Team %s not found", name)
	}

	return team
}

func Compare(t *testing.T, expected, value interface{}) {
	err := utils.Compare(expected, value)
	if err != nil {
		Fatalf(t, "Failed to compare values: %v", err)
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
		Fatalf(t, "Failed to get current directory: %v", err)
	}

	return filepath.Base(dir)
}

func CreateDir(t *testing.T, dir string) {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		Fatalf(t, "Failed to create directory %s: %v", dir, err)
	}
}

func CreateFile(t *testing.T, file string, content string) {
	f, err := os.Create(file)
	if err != nil {
		Fatalf(t, "Failed to create file %s: %v", file, err)
	}
	defer func() {
		err := f.Close()
		if err != nil {
			t.Fatalf("Failed to close file %s: %v", file, err)
		}
	}()

	_, err = f.WriteString(content)
	if err != nil {
		Fatalf(t, "Failed to write content to file %s: %v", file, err)
	}
}

func HashFile(t *testing.T, file string) string {
	f, err := os.Open(file)
	if err != nil {
		Fatalf(t, "Failed to open file %s: %v", file, err)
	}
	defer func() {
		err := f.Close()
		if err != nil {
			t.Fatalf("Failed to close file %s: %v", file, err)
		}
	}()

	hash, err := crypto_utils.HashFile(f)
	if err != nil {
		Fatalf(t, "Failed to hash file %s: %v", file, err)
	}

	return hash
}

func Format(msg string, a ...any) string {
	res := msg[:]
	for i, v := range a {
		placeholder := fmt.Sprintf("{%d}", i)
		res = strings.ReplaceAll(res, placeholder, fmt.Sprint(v))
	}
	return res
}
