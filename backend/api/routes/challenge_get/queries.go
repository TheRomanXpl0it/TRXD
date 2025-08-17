package challenge_get

import (
	"context"
	"database/sql"
	"strings"
	"trxd/db"
	"trxd/db/sqlc"
)

type Chall struct {
	ID          int32                         `json:"id"`
	Name        string                        `json:"name"`
	Category    string                        `json:"category"`
	Description string                        `json:"description"`
	Difficulty  string                        `json:"difficulty"`
	Authors     []string                      `json:"authors"`
	Instance    bool                          `json:"instance"`
	Hidden      bool                          `json:"hidden"`
	Points      int32                         `json:"points"`
	Solves      int32                         `json:"solves"`
	FirstBlood  *sqlc.GetFirstBloodRow        `json:"first_blood"`
	Host        string                        `json:"host"`
	Port        int32                         `json:"port"`
	Attachments []string                      `json:"attachments"`
	Tags        []string                      `json:"tags"`
	Flags       []sqlc.GetFlagsByChallengeRow `json:"flags"`
	Timeout     int32                         `json:"timeout"`
	Solved      bool                          `json:"solved"`
	SolvesList  []sqlc.GetChallengeSolvesRow  `json:"solves_list"`
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

	flags, err := GetFlagsByChallenge(ctx, challenge.ID)
	if err != nil {
		return nil, err
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
		chall.Authors = strings.Split(challenge.Authors.String, ",") // TODO: change char
	}
	chall.Instance = challenge.Type != sqlc.DeployTypeNormal
	chall.Hidden = challenge.Hidden
	chall.Points = challenge.Points
	chall.Solves = challenge.Solves
	if challenge.Host.Valid {
		chall.Host = challenge.Host.String
	}
	if challenge.Port.Valid {
		chall.Port = challenge.Port.Int32
	}
	chall.Attachments = []string{}
	if challenge.Attachments.Valid {
		chall.Attachments = strings.Split(challenge.Attachments.String, ",") // TODO: change char
	}
	if author {
		if flags != nil {
			chall.Flags = flags
		} else {
			chall.Flags = []sqlc.GetFlagsByChallengeRow{}
		}
	}
	chall.Tags = []string{}
	if tags != nil {
		chall.Tags = tags
	}
	chall.Solved = solved

	solves, err := db.Sql.GetChallengeSolves(ctx, id)
	if err != nil {
		return nil, err
	}
	if solves == nil {
		solves = []sqlc.GetChallengeSolvesRow{}
	}
	chall.SolvesList = solves

	firstBlood, err := GetFirstBlood(ctx, id)
	if err != nil {
		return nil, err
	}
	chall.FirstBlood = firstBlood

	return &chall, nil
}
