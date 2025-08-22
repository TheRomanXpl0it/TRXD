package challenges_get

import (
	"context"
	"database/sql"
	"strings"
	"trxd/db"
	"trxd/db/sqlc"
	"trxd/utils/consts"
)

type DockerConfig struct {
	Image      string `json:"image,omitempty"`
	Compose    string `json:"compose,omitempty"`
	HashDomain *bool  `json:"hash_domain,omitempty"`
	Lifetime   *int   `json:"lifetime,omitempty"`
	Envs       string `json:"envs,omitempty"`
	MaxMemory  *int   `json:"max_memory,omitempty"`
	MaxCpu     string `json:"max_cpu,omitempty"`
}

type Chall struct {
	ID           int32                          `json:"id"`
	Name         string                         `json:"name"`
	Category     string                         `json:"category"`
	Description  string                         `json:"description"`
	Difficulty   string                         `json:"difficulty"`
	Authors      []string                       `json:"authors"`
	Instance     bool                           `json:"instance"`
	Type         *sqlc.DeployType               `json:"type,omitempty"`
	Hidden       *bool                          `json:"hidden,omitempty"`
	MaxPoints    *int                           `json:"max_points,omitempty"`
	ScoreType    *sqlc.ScoreType                `json:"score_type,omitempty"`
	Points       int                            `json:"points"`
	Solves       int                            `json:"solves"`
	FirstBlood   *sqlc.GetFirstBloodRow         `json:"first_blood"`
	Host         string                         `json:"host"`
	Port         int                            `json:"port"`
	Attachments  []string                       `json:"attachments"`
	Tags         []string                       `json:"tags"`
	Flags        *[]sqlc.GetFlagsByChallengeRow `json:"flags,omitempty"`
	Timeout      int                            `json:"timeout"`
	Solved       bool                           `json:"solved"`
	SolvesList   []sqlc.GetChallengeSolvesRow   `json:"solves_list"`
	DockerConfig *DockerConfig                  `json:"docker_config,omitempty"`
}

func GetFlagsByChallenge(ctx context.Context, challengeID int32) ([]sqlc.GetFlagsByChallengeRow, error) {
	flags, err := db.Sql.GetFlagsByChallenge(ctx, challengeID)
	if err != nil {
		return nil, err
	}

	return flags, nil
}

func IsChallengeSolved(ctx context.Context, id int32, uid int32) (bool, error) {
	solved, err := db.Sql.IsChallengeSolved(ctx, sqlc.IsChallengeSolvedParams{
		ChallID: id,
		ID:      uid,
	})
	if err != nil {
		return false, err
	}

	return solved, nil
}

func GetFirstBlood(ctx context.Context, id int32) (*sqlc.GetFirstBloodRow, error) {
	firstBlood, err := db.Sql.GetFirstBlood(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &firstBlood, nil
}

func GetChallenge(ctx context.Context, id int32, uid int32, author bool) (*Chall, error) {
	var chall Chall

	challenge, err := db.GetChallengeByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if challenge == nil {
		return nil, nil
	}
	if !author && challenge.Hidden {
		return nil, nil
	}

	tags, err := db.GetTagsByChallenge(ctx, challenge.ID)
	if err != nil {
		return nil, err
	}

	solved, err := IsChallengeSolved(ctx, id, uid)
	if err != nil {
		return nil, err
	}

	chall.ID = challenge.ID
	chall.Name = challenge.Name
	chall.Category = challenge.Category
	chall.Description = challenge.Description
	if challenge.Difficulty.Valid {
		chall.Difficulty = challenge.Difficulty.String
	}
	chall.Authors = []string{}
	if challenge.Authors.Valid {
		chall.Authors = strings.Split(challenge.Authors.String, consts.Separator)
	}
	chall.Instance = challenge.Type != sqlc.DeployTypeNormal
	chall.Points = int(challenge.Points)
	chall.Solves = int(challenge.Solves)
	if challenge.Host.Valid {
		chall.Host = challenge.Host.String // TODO: override with instance
	}
	if challenge.Port.Valid {
		chall.Port = int(challenge.Port.Int32) // TODO: override with instance
	}
	chall.Attachments = []string{}
	if challenge.Attachments.Valid {
		chall.Attachments = strings.Split(challenge.Attachments.String, consts.Separator)
	}
	chall.Tags = []string{}
	if tags != nil {
		chall.Tags = tags
	}
	chall.Solved = solved
	// TODO: add chall.Timeout from the instance

	solves, err := db.Sql.GetChallengeSolves(ctx, id)
	if err != nil {
		return nil, err
	}
	chall.SolvesList = []sqlc.GetChallengeSolvesRow{}
	if solves != nil {
		chall.SolvesList = solves
	}

	firstBlood, err := GetFirstBlood(ctx, id)
	if err != nil {
		return nil, err
	}
	chall.FirstBlood = firstBlood

	if !author {
		return &chall, nil
	}

	flags, err := GetFlagsByChallenge(ctx, challenge.ID)
	if err != nil {
		return nil, err
	}

	maxPoints := int(challenge.MaxPoints)
	chall.Hidden = &challenge.Hidden
	chall.Type = &challenge.Type
	chall.MaxPoints = &maxPoints
	chall.ScoreType = &challenge.ScoreType
	chall.Flags = &[]sqlc.GetFlagsByChallengeRow{}
	if flags != nil {
		chall.Flags = &flags
	}

	dockerConfig, err := db.Sql.GetChallDockerConfig(ctx, challenge.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &chall, nil
		}
		return nil, err
	}

	chall.DockerConfig = &DockerConfig{}

	if dockerConfig.Image.Valid {
		chall.DockerConfig.Image = dockerConfig.Image.String
	}
	if dockerConfig.Compose.Valid {
		chall.DockerConfig.Compose = dockerConfig.Compose.String
	}
	chall.DockerConfig.HashDomain = &dockerConfig.HashDomain
	if dockerConfig.Lifetime.Valid {
		lifetime := int(dockerConfig.Lifetime.Int32)
		chall.DockerConfig.Lifetime = &lifetime
	}
	if dockerConfig.Envs.Valid {
		chall.DockerConfig.Envs = dockerConfig.Envs.String
	}
	if dockerConfig.MaxMemory.Valid {
		maxMemory := int(dockerConfig.MaxMemory.Int32)
		chall.DockerConfig.MaxMemory = &maxMemory
	}
	if dockerConfig.MaxCpu.Valid {
		chall.DockerConfig.MaxCpu = dockerConfig.MaxCpu.String
	}

	return &chall, nil
}
