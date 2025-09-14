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
	Image      string  `json:"image"`
	Compose    string  `json:"compose"`
	HashDomain *bool   `json:"hash_domain"`
	Lifetime   *int    `json:"lifetime"`
	Envs       *string `json:"envs"`
	MaxMemory  *int    `json:"max_memory"`
	MaxCpu     *string `json:"max_cpu"`
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
	FirstBlood  bool                         `json:"first_blood"`
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
		Difficulty:  challenge.Difficulty,
		Authors:     []string{},
		Instance:    challenge.Type != sqlc.DeployTypeNormal,
		Points:      int(challenge.Points),
		Solves:      int(challenge.Solves),
		Attachments: []string{},
		Tags:        []string{},
		Solved:      solved,
		SolvesList:  []sqlc.GetChallengeSolvesRow{},
		Host:        challenge.Host,
		Port:        int(challenge.Port),
	}

	if instance != nil {
		chall.Host = instance.Host
		if instance.Port.Valid {
			chall.Port = int(instance.Port.Int32)
		}
		chall.Timeout = int(time.Until(instance.ExpiresAt).Seconds())
		if chall.Timeout < 0 {
			chall.Timeout = 0
		}
	}
	if challenge.Authors != "" {
		chall.Authors = strings.Split(challenge.Authors, consts.Separator)
	}
	if challenge.Attachments != "" {
		chall.Attachments = strings.Split(challenge.Attachments, consts.Separator)
	}
	if tags != nil {
		chall.Tags = tags
	}
	if solves != nil {
		chall.SolvesList = solves
		chall.FirstBlood = solves[0].ID == tid
	}

	var noDockerConfig bool
	dockerConfig, err := db.Sql.GetChallDockerConfig(ctx, challenge.ID)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		} else {
			noDockerConfig = err == sql.ErrNoRows
		}
	}

	if dockerConfig.HashDomain {
		chall.Port = 0
	}

	if !author { // Not Author
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

	if noDockerConfig {
		return &chall, nil
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
