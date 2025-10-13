package challenges_all_get

import (
	"context"
	"database/sql"
	"strings"
	"time"
	"trxd/db"
	"trxd/db/sqlc"
	"trxd/utils/consts"
)

type Chall struct {
	ID          int32          `json:"id"`
	Name        string         `json:"name"`
	Category    string         `json:"category"`
	Description string         `json:"description"`
	Difficulty  string         `json:"difficulty"`
	Authors     []string       `json:"authors"`
	Instance    bool           `json:"instance"`
	Hidden      bool           `json:"hidden"`
	Points      int            `json:"points"`
	Solves      int            `json:"solves"`
	FirstBlood  bool           `json:"first_blood"`
	Host        string         `json:"host"`
	Port        int            `json:"port"`
	Attachments []string       `json:"attachments"`
	Tags        []string       `json:"tags"`
	Solved      bool           `json:"solved"`
	MaxPoints   int            `json:"max_points"`
	ScoreType   sqlc.ScoreType `json:"score_type"`
	Timeout     int            `json:"timeout"`
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

func GetChallenge(ctx context.Context, challenge *sqlc.Challenge, uid int32, tid int32) (*Chall, error) {
	tags, err := db.GetTagsByChallenge(ctx, challenge.ID)
	if err != nil {
		return nil, err
	}

	solved, err := IsChallengeSolved(ctx, challenge.ID, uid)
	if err != nil {
		return nil, err
	}

	instance, err := GetInstanceInfo(ctx, challenge.ID, tid)
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
		Hidden:      challenge.Hidden,
		Points:      int(challenge.Points),
		Solves:      int(challenge.Solves),
		// TODO: first blood
		Attachments: []string{},
		Tags:        []string{},
		Solved:      solved,
		Host:        challenge.Host,
		Port:        int(challenge.Port),
		MaxPoints:   int(challenge.MaxPoints),
		ScoreType:   challenge.ScoreType,
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

	return &chall, nil
}

func GetChallenges(ctx context.Context, uid int32, tid int32, author bool) ([]*Chall, error) {
	challenges, err := db.Sql.GetChallenges(ctx)
	if err != nil {
		return nil, err
	}

	challsData := make([]*Chall, 0)
	for _, challenge := range challenges {
		if !author && challenge.Hidden {
			continue
		}

		chall, err := GetChallenge(ctx, &challenge, uid, tid)
		if err != nil {
			return nil, err
		}

		challsData = append(challsData, chall)
	}

	return challsData, nil
}
