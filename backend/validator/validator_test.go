package validator_test

import (
	"fmt"
	"math"
	"strings"
	"testing"
	"trxd/db/sqlc"
	"trxd/utils/consts"
	"trxd/utils/test_utils"
	"trxd/validator"
)

func varTest(t *testing.T, tag string, v interface{}, results ...string) {
	var result string
	if len(results) != 0 {
		result = results[0]
	}

	ok, err := validator.Var(nil, v, tag)
	if err != nil {
		if err.Error() != result {
			t.Fatalf("Validator failed on: %s = <%v>: %v != %v", tag, v, result, err)
		}
	} else if !ok {
		t.Fatalf("Validator returned wrong result on: %s = <%v>: got <%v>, want %v", tag, v, ok, result)
	} else if result != "" {
		t.Fatalf("Validator should have failed on: %s = <%v> (%s)", tag, v, result)
	}
}

func TestValidators(t *testing.T) {
	varTest(t, "id", -1, test_utils.Format(consts.MinError, "id", 0))
	varTest(t, "id", 0)
	varTest(t, "id", 1337)
	varTest(t, "id", math.MaxInt32)
	varTest(t, "id", math.MaxInt32+1, test_utils.Format(consts.MaxError, "id", math.MaxInt32))

	varTest(t, "password", "", test_utils.Format(consts.MinError, "password", consts.MinPasswordLen))
	varTest(t, "password", strings.Repeat("a", consts.MinPasswordLen-1), test_utils.Format(consts.MinError, "password", consts.MinPasswordLen))
	varTest(t, "password", strings.Repeat("a", consts.MinPasswordLen))
	varTest(t, "password", strings.Repeat("a", consts.MaxPasswordLen))
	varTest(t, "password", strings.Repeat("a", consts.MaxPasswordLen+1), test_utils.Format(consts.MaxError, "password", consts.MaxPasswordLen))

	varTest(t, "country", "")
	varTest(t, "country", "a", consts.InvalidCountry)
	varTest(t, "country", "aaa", consts.InvalidCountry)
	varTest(t, "country", "USA")
	varTest(t, "country", "JPN")
	varTest(t, "country", "aaaa", consts.InvalidCountry)

	varTest(t, "category_name", "")
	varTest(t, "category_name", "a")
	varTest(t, "category_name", strings.Repeat("a", consts.MaxCategoryLen))
	varTest(t, "category_name", strings.Repeat("a", consts.MaxCategoryLen+1), test_utils.Format(consts.MaxError, "category_name", consts.MaxCategoryLen))

	varTest(t, "challenge_name", "")
	varTest(t, "challenge_name", "a")
	varTest(t, "challenge_name", strings.Repeat("a", consts.MaxChallNameLen))
	varTest(t, "challenge_name", strings.Repeat("a", consts.MaxChallNameLen+1), test_utils.Format(consts.MaxError, "challenge_name", consts.MaxChallNameLen))

	varTest(t, "challenge_description", "")
	varTest(t, "challenge_description", "a")
	varTest(t, "challenge_description", strings.Repeat("a", consts.MaxChallDescLen))
	varTest(t, "challenge_description", strings.Repeat("a", consts.MaxChallDescLen+1), test_utils.Format(consts.MaxError, "challenge_description", consts.MaxChallDescLen))

	varTest(t, "challenge_difficulty", "")
	varTest(t, "challenge_difficulty", "a")
	varTest(t, "challenge_difficulty", strings.Repeat("a", consts.MaxChallDifficultyLen))
	varTest(t, "challenge_difficulty", strings.Repeat("a", consts.MaxChallDifficultyLen+1), test_utils.Format(consts.MaxError, "challenge_difficulty", consts.MaxChallDifficultyLen))

	varTest(t, "challenge_authors", []string{})
	varTest(t, "challenge_authors", []string{""})
	varTest(t, "challenge_authors", []string{"a"})
	varTest(t, "challenge_authors", []string{"a", ""})
	varTest(t, "challenge_authors", []string{"a", "a"})
	varTest(t, "challenge_authors", []string{strings.Repeat("a", consts.MaxAuthorNameLen)})
	varTest(t, "challenge_authors", []string{strings.Repeat("a", consts.MaxAuthorNameLen+1)}, test_utils.Format(consts.MaxError, "[0]", consts.MaxAuthorNameLen))

	varTest(t, "challenge_tags", []string{})
	varTest(t, "challenge_tags", []string{""})
	varTest(t, "challenge_tags", []string{"a"})
	varTest(t, "challenge_tags", []string{"a", ""})
	varTest(t, "challenge_tags", []string{"a", "a"})
	varTest(t, "challenge_tags", []string{strings.Repeat("a", consts.MaxTagNameLen)})
	varTest(t, "challenge_tags", []string{strings.Repeat("a", consts.MaxTagNameLen+1)}, test_utils.Format(consts.MaxError, "[0]", consts.MaxTagNameLen))

	varTest(t, "challenge_type", "", test_utils.Format(consts.OneOfError, "challenge_type", strings.Join(consts.DeployTypesStr, " ")))
	varTest(t, "challenge_type", sqlc.DeployTypeNormal)
	varTest(t, "challenge_type", sqlc.DeployTypeContainer)
	varTest(t, "challenge_type", sqlc.DeployTypeCompose)
	varTest(t, "challenge_type", "aaa", test_utils.Format(consts.OneOfError, "challenge_type", strings.Join(consts.DeployTypesStr, " ")))

	varTest(t, "challenge_max_points", -1, test_utils.Format(consts.MinError, "challenge_max_points", 0))
	varTest(t, "challenge_max_points", 0)
	varTest(t, "challenge_max_points", 1337)
	varTest(t, "challenge_max_points", math.MaxInt32)
	varTest(t, "challenge_max_points", math.MaxInt32+1, test_utils.Format(consts.MaxError, "challenge_max_points", math.MaxInt32))

	varTest(t, "challenge_score_type", "", test_utils.Format(consts.OneOfError, "challenge_score_type", strings.Join(consts.ScoreTypesStr, " ")))
	varTest(t, "challenge_score_type", sqlc.ScoreTypeStatic)
	varTest(t, "challenge_score_type", sqlc.ScoreTypeDynamic)
	varTest(t, "challenge_score_type", "aaa", test_utils.Format(consts.OneOfError, "challenge_score_type", strings.Join(consts.ScoreTypesStr, " ")))

	varTest(t, "challenge_port", consts.MinPort-1, test_utils.Format(consts.MinError, "challenge_port", consts.MinPort))
	varTest(t, "challenge_port", consts.MinPort)
	varTest(t, "challenge_port", consts.MaxPort)
	varTest(t, "challenge_port", consts.MaxPort+1, test_utils.Format(consts.MaxError, "challenge_port", consts.MaxPort))

	varTest(t, "challenge_lifetime", -1, test_utils.Format(consts.MinError, "challenge_lifetime", 0))
	varTest(t, "challenge_lifetime", 0)
	varTest(t, "challenge_lifetime", 1337)
	varTest(t, "challenge_lifetime", math.MaxInt32)
	varTest(t, "challenge_lifetime", math.MaxInt32+1, test_utils.Format(consts.MaxError, "challenge_lifetime", math.MaxInt32))

	varTest(t, "challenge_envs", "")
	varTest(t, "challenge_envs", "a", consts.InvalidEnvs)
	varTest(t, "challenge_envs", "[]")
	varTest(t, "challenge_envs", "{}")
	varTest(t, "challenge_envs", `{"key":"value","key2":"value2"}`)
	varTest(t, "challenge_envs", `{"key": "value", "key2": "value2"}`)

	varTest(t, "challenge_max_memory", -1, test_utils.Format(consts.MinError, "challenge_max_memory", 0))
	varTest(t, "challenge_max_memory", 0)
	varTest(t, "challenge_max_memory", 1337)
	varTest(t, "challenge_max_memory", math.MaxInt32)
	varTest(t, "challenge_max_memory", math.MaxInt32+1, test_utils.Format(consts.MaxError, "challenge_max_memory", math.MaxInt32))

	varTest(t, "challenge_max_cpu", "-1", consts.InvalidMaxCpu)
	varTest(t, "challenge_max_cpu", "0", consts.InvalidMaxCpu)
	varTest(t, "challenge_max_cpu", "13.37")
	varTest(t, "challenge_max_cpu", "1337")
	varTest(t, "challenge_max_cpu", fmt.Sprint(math.MaxInt32))
	varTest(t, "challenge_max_cpu", fmt.Sprint(math.MaxInt32+1), consts.InvalidMaxCpu)

	varTest(t, "attachments", []string{})
	varTest(t, "attachments", []string{""})
	varTest(t, "attachments", []string{"a"})
	varTest(t, "attachments", []string{"a", ""})
	varTest(t, "attachments", []string{"a", "a"})
	varTest(t, "attachments", []string{strings.Repeat("a", consts.MaxAttachmentNameLen)})
	varTest(t, "attachments", []string{strings.Repeat("a", consts.MaxAttachmentNameLen+1)}, test_utils.Format(consts.MaxError, "[0]", consts.MaxAttachmentNameLen))

	varTest(t, "flag", "")
	varTest(t, "flag", "a")
	varTest(t, "flag", strings.Repeat("a", consts.MaxFlagLen))
	varTest(t, "flag", strings.Repeat("a", consts.MaxFlagLen+1), test_utils.Format(consts.MaxError, "flag", consts.MaxFlagLen))

	varTest(t, "team_name", "")
	varTest(t, "team_name", "a")
	varTest(t, "team_name", strings.Repeat("a", consts.MaxTeamNameLen))
	varTest(t, "team_name", strings.Repeat("a", consts.MaxTeamNameLen+1), test_utils.Format(consts.MaxError, "team_name", consts.MaxTeamNameLen))

	varTest(t, "user_name", "")
	varTest(t, "user_name", "a")
	varTest(t, "user_name", strings.Repeat("a", consts.MaxUserNameLen))
	varTest(t, "user_name", strings.Repeat("a", consts.MaxUserNameLen+1), test_utils.Format(consts.MaxError, "user_name", consts.MaxUserNameLen))

	varTest(t, "user_email", "", consts.InvalidEmail)
	varTest(t, "user_email", "a", consts.InvalidEmail)
	varTest(t, "user_email", "test@example.com")
	varTest(t, "user_email", "test+alias@example.com")
	varTest(t, "user_email", strings.Repeat("a", consts.MaxEmailLen+1), test_utils.Format(consts.MaxError, "user_email", consts.MaxEmailLen))
}
