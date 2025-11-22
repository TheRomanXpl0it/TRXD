package validator_test

import (
	"fmt"
	"math"
	"strings"
	"testing"
	"trxd/api/validator"
	"trxd/db/sqlc"
	"trxd/utils/consts"
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
	varTest(t, "id", -1, "id must be at least 0")
	varTest(t, "id", 0)
	varTest(t, "id", 1337)
	varTest(t, "id", math.MaxInt32)
	varTest(t, "id", math.MaxInt32+1, "id must not exceed 2147483647")

	varTest(t, "password", "", "password must be at least 8")
	varTest(t, "password", strings.Repeat("a", consts.MinPasswordLen-1), "password must be at least 8")
	varTest(t, "password", strings.Repeat("a", consts.MinPasswordLen))
	varTest(t, "password", strings.Repeat("a", consts.MaxPasswordLen))
	varTest(t, "password", strings.Repeat("a", consts.MaxPasswordLen+1), "password must not exceed 64")

	varTest(t, "country", "")
	varTest(t, "country", "a", consts.InvalidCountry)
	varTest(t, "country", "aaa", consts.InvalidCountry)
	varTest(t, "country", "USA")
	varTest(t, "country", "JPN")
	varTest(t, "country", "aaaa", consts.InvalidCountry)

	varTest(t, "valid_http_url", "")
	varTest(t, "valid_http_url", "a", consts.InvalidHttpUrl)
	varTest(t, "valid_http_url", "file://a", consts.InvalidHttpUrl)
	varTest(t, "valid_http_url", "http://a")
	varTest(t, "valid_http_url", "https://a")
	varTest(t, "valid_http_url", "http://example.com")
	varTest(t, "valid_http_url", "https://example.com")

	varTest(t, "image_url", "")
	varTest(t, "image_url", "a", consts.InvalidHttpUrl)
	varTest(t, "image_url", "file://a", consts.InvalidHttpUrl)
	varTest(t, "image_url", "http://a")
	varTest(t, "image_url", "https://a")
	varTest(t, "image_url", "http://example.com")
	varTest(t, "image_url", "https://example.com")
	varTest(t, "image_url", strings.Repeat("a", consts.MaxImageLen+1), "image_url must not exceed 1024")

	varTest(t, "category_name", "")
	varTest(t, "category_name", "a")
	varTest(t, "category_name", strings.Repeat("a", consts.MaxCategoryLen))
	varTest(t, "category_name", strings.Repeat("a", consts.MaxCategoryLen+1), "category_name must not exceed 32")

	varTest(t, "category_icon", "")
	varTest(t, "category_icon", "a")
	varTest(t, "category_icon", strings.Repeat("a", consts.MaxIconLen))
	varTest(t, "category_icon", strings.Repeat("a", consts.MaxIconLen+1), "category_icon must not exceed 32")

	varTest(t, "challenge_name", "")
	varTest(t, "challenge_name", "a")
	varTest(t, "challenge_name", strings.Repeat("a", consts.MaxChallNameLen))
	varTest(t, "challenge_name", strings.Repeat("a", consts.MaxChallNameLen+1), "challenge_name must not exceed 128")

	varTest(t, "challenge_description", "")
	varTest(t, "challenge_description", "a")
	varTest(t, "challenge_description", strings.Repeat("a", consts.MaxChallDescLen))
	varTest(t, "challenge_description", strings.Repeat("a", consts.MaxChallDescLen+1), "challenge_description must not exceed 1024")

	varTest(t, "challenge_difficulty", "")
	varTest(t, "challenge_difficulty", "a")
	varTest(t, "challenge_difficulty", strings.Repeat("a", consts.MaxChallDifficultyLen))
	varTest(t, "challenge_difficulty", strings.Repeat("a", consts.MaxChallDifficultyLen+1), "challenge_difficulty must not exceed 16")

	varTest(t, "challenge_type", "", "challenge_type must be one of: Normal Container Compose")
	varTest(t, "challenge_type", sqlc.DeployTypeNormal)
	varTest(t, "challenge_type", sqlc.DeployTypeContainer)
	varTest(t, "challenge_type", sqlc.DeployTypeCompose)
	varTest(t, "challenge_type", "aaa", "challenge_type must be one of: Normal Container Compose")

	varTest(t, "challenge_max_points", -1, "challenge_max_points must be at least 0")
	varTest(t, "challenge_max_points", 0)
	varTest(t, "challenge_max_points", 1337)
	varTest(t, "challenge_max_points", math.MaxInt32)
	varTest(t, "challenge_max_points", math.MaxInt32+1, "challenge_max_points must not exceed 2147483647")

	varTest(t, "challenge_score_type", "", "challenge_score_type must be one of: Static Dynamic")
	varTest(t, "challenge_score_type", sqlc.ScoreTypeStatic)
	varTest(t, "challenge_score_type", sqlc.ScoreTypeDynamic)
	varTest(t, "challenge_score_type", "aaa", "challenge_score_type must be one of: Static Dynamic")

	varTest(t, "challenge_port", consts.MinPort-1, "challenge_port must be at least 0")
	varTest(t, "challenge_port", consts.MinPort)
	varTest(t, "challenge_port", consts.MaxPort)
	varTest(t, "challenge_port", consts.MaxPort+1, "challenge_port must not exceed 65535")

	varTest(t, "challenge_lifetime", -1, "challenge_lifetime must be at least 0")
	varTest(t, "challenge_lifetime", 0)
	varTest(t, "challenge_lifetime", 1337)
	varTest(t, "challenge_lifetime", math.MaxInt32)
	varTest(t, "challenge_lifetime", math.MaxInt32+1, "challenge_lifetime must not exceed 2147483647")

	varTest(t, "challenge_envs", "")
	varTest(t, "challenge_envs", "a", consts.InvalidEnvs)
	varTest(t, "challenge_envs", "[]")
	varTest(t, "challenge_envs", "{}")
	varTest(t, "challenge_envs", `{"key":"value","key2":"value2"}`)
	varTest(t, "challenge_envs", `{"key": "value", "key2": "value2"}`)

	varTest(t, "challenge_max_memory", -1, "challenge_max_memory must be at least 0")
	varTest(t, "challenge_max_memory", 0)
	varTest(t, "challenge_max_memory", 1337)
	varTest(t, "challenge_max_memory", math.MaxInt32)
	varTest(t, "challenge_max_memory", math.MaxInt32+1, "challenge_max_memory must not exceed 2147483647")

	varTest(t, "challenge_max_cpu", "-1", consts.InvalidMaxCpu)
	varTest(t, "challenge_max_cpu", "0", consts.InvalidMaxCpu)
	varTest(t, "challenge_max_cpu", "13.37")
	varTest(t, "challenge_max_cpu", "1337")
	varTest(t, "challenge_max_cpu", fmt.Sprint(math.MaxInt32))
	varTest(t, "challenge_max_cpu", fmt.Sprint(math.MaxInt32+1), consts.InvalidMaxCpu)

	varTest(t, "flag", "")
	varTest(t, "flag", "a")
	varTest(t, "flag", strings.Repeat("a", consts.MaxFlagLen))
	varTest(t, "flag", strings.Repeat("a", consts.MaxFlagLen+1), "flag must not exceed 128")

	varTest(t, "tag_name", "")
	varTest(t, "tag_name", "a")
	varTest(t, "tag_name", strings.Repeat("a", consts.MaxTagNameLen))
	varTest(t, "tag_name", strings.Repeat("a", consts.MaxTagNameLen+1), "tag_name must not exceed 32")

	varTest(t, "team_name", "")
	varTest(t, "team_name", "a")
	varTest(t, "team_name", strings.Repeat("a", consts.MaxTeamNameLen))
	varTest(t, "team_name", strings.Repeat("a", consts.MaxTeamNameLen+1), "team_name must not exceed 64")

	varTest(t, "team_bio", "")
	varTest(t, "team_bio", "a")
	varTest(t, "team_bio", strings.Repeat("a", consts.MaxBioLen))
	varTest(t, "team_bio", strings.Repeat("a", consts.MaxBioLen+1), "team_bio must not exceed 10240")

	varTest(t, "user_name", "")
	varTest(t, "user_name", "a")
	varTest(t, "user_name", strings.Repeat("a", consts.MaxUserNameLen))
	varTest(t, "user_name", strings.Repeat("a", consts.MaxUserNameLen+1), "user_name must not exceed 64")

	varTest(t, "user_email", "", consts.InvalidEmail)
	varTest(t, "user_email", "a", consts.InvalidEmail)
	varTest(t, "user_email", "test@example.com")
	varTest(t, "user_email", "test+alias@example.com")
	varTest(t, "user_email", strings.Repeat("a", consts.MaxEmailLen+1), "user_email must not exceed 256")
}
