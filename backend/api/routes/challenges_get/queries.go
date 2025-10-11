package challenges_get

import (
	"context"
	"database/sql"
	"trxd/db"
	"trxd/db/sqlc"
)

type DockerConfig struct {
	Image      string  `json:"image"`
	Compose    string  `json:"compose"`
	HashDomain *bool   `json:"hash_domain"`
	Lifetime   *int    `json:"lifetime"`
	Envs       *string `json:"envs"`
	MaxMemory  *int    `json:"max_memory"`
	MaxCpu     *string `json:"max_cpu"`
}

type Chall struct {
	SolvesList []sqlc.GetChallengeSolvesRow `json:"solves_list"`

	Type         *sqlc.DeployType               `json:"type,omitempty"`
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

	solves, err := db.Sql.GetChallengeSolves(ctx, id)
	if err != nil {
		return nil, err
	}

	chall := Chall{
		SolvesList: []sqlc.GetChallengeSolvesRow{},
	}
	if solves != nil {
		chall.SolvesList = solves
	}

	if !author { // Not Author
		return &chall, nil
	}

	flags, err := GetFlagsByChallenge(ctx, challenge.ID)
	if err != nil {
		return nil, err
	}

	chall.Type = &challenge.Type
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

	lifetime := int(dockerConfig.Lifetime)
	maxMemory := int(dockerConfig.MaxMemory)
	chall.DockerConfig = &DockerConfig{
		Image:      dockerConfig.Image,
		Compose:    dockerConfig.Compose,
		HashDomain: &dockerConfig.HashDomain,
		Lifetime:   &lifetime,
		Envs:       &dockerConfig.Envs,
		MaxMemory:  &maxMemory,
		MaxCpu:     &dockerConfig.MaxCpu,
	}

	return &chall, nil
}
