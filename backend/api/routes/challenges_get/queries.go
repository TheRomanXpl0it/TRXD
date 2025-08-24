package challenges_get

import (
	"context"
	"database/sql"
	"strings"
	"time"
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
	ID          int32                        `json:"id"`
	Name        string                       `json:"name"`
	Category    string                       `json:"category"`
	Description string                       `json:"description"`
	Difficulty  string                       `json:"difficulty"`
	Authors     []string                     `json:"authors"`
	Instance    bool                         `json:"instance"`
	Points      int                          `json:"points"`
	Solves      int                          `json:"solves"`
	FirstBlood  *sqlc.GetFirstBloodRow       `json:"first_blood"`
	Host        string                       `json:"host"`
	Port        int                          `json:"port"`
	Attachments []string                     `json:"attachments"`
	Tags        []string                     `json:"tags"`
	Timeout     int                          `json:"timeout"`
	Solved      bool                         `json:"solved"`
	SolvesList  []sqlc.GetChallengeSolvesRow `json:"solves_list"`

	Type         *sqlc.DeployType               `json:"type,omitempty"`
	Hidden       *bool                          `json:"hidden,omitempty"`
	MaxPoints    *int                           `json:"max_points,omitempty"`
	ScoreType    *sqlc.ScoreType                `json:"score_type,omitempty"`
	Flags        *[]sqlc.GetFlagsByChallengeRow `json:"flags,omitempty"`
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

func GetInstanceInfo(ctx context.Context, challID int32, teamID int32) (*sqlc.GetInstanceInfoRow, error) {
	instance, err := db.Sql.GetInstanceInfo(ctx, sqlc.GetInstanceInfoParams{
		TeamID:  teamID,
		ChallID: challID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &instance, nil
}

func GetChallenge(ctx context.Context, id int32, uid int32, tid int32, author bool) (*Chall, error) {
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

	firstBlood, err := GetFirstBlood(ctx, id)
	if err != nil {
		return nil, err
	}

	solves, err := db.Sql.GetChallengeSolves(ctx, id)
	if err != nil {
		return nil, err
	}

	instance, err := GetInstanceInfo(ctx, id, tid)
	if err != nil {
		return nil, err
	}

	chall := Chall{
		ID:          challenge.ID,
		Name:        challenge.Name,
		Category:    challenge.Category,
		Description: challenge.Description,
		Authors:     []string{},
		Instance:    challenge.Type != sqlc.DeployTypeNormal,
		Points:      int(challenge.Points),
		Solves:      int(challenge.Solves),
		Attachments: []string{},
		Tags:        []string{},
		Solved:      solved,
		SolvesList:  []sqlc.GetChallengeSolvesRow{},
		FirstBlood:  firstBlood,
	}

	if challenge.Difficulty.Valid {
		chall.Difficulty = challenge.Difficulty.String
	}
	if challenge.Authors.Valid {
		chall.Authors = strings.Split(challenge.Authors.String, consts.Separator)
	}
	if instance != nil {
		chall.Host = instance.Host
	} else if challenge.Host.Valid {
		chall.Host = challenge.Host.String
	}
	if instance != nil {
		chall.Port = int(instance.Port)
	} else if challenge.Port.Valid {
		chall.Port = int(challenge.Port.Int32)
	}
	if challenge.Attachments.Valid {
		chall.Attachments = strings.Split(challenge.Attachments.String, consts.Separator)
	}
	if tags != nil {
		chall.Tags = tags
	}
	if instance != nil {
		chall.Timeout = int(time.Until(instance.ExpiresAt).Seconds())
		if chall.Timeout < 0 {
			chall.Timeout = 0
		}
	}
	if solves != nil {
		chall.SolvesList = solves
	}

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
