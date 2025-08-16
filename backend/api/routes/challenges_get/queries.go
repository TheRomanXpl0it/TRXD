package challenges_get

import (
	"context"
	"trxd/db"
	"trxd/db/sqlc"
)

type Chall struct {
	ID         int32    `json:"id"`
	Name       string   `json:"name"`
	Category   string   `json:"category"`
	Difficulty string   `json:"difficulty"`
	Instance   bool     `json:"instance"`
	Hidden     bool     `json:"hidden"`
	Points     int32    `json:"points"`
	Solves     int32    `json:"solves"`
	Tags       []string `json:"tags"`
	Solved     bool     `json:"solved"`
}

func GetChallengePreview(ctx context.Context, id int32, uid int32, author bool) (*Chall, error) {
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

	solved, err := db.IsChallengeSolved(ctx, id, uid)
	if err != nil {
		return nil, err
	}

	chall.ID = challenge.ID
	chall.Name = challenge.Name
	chall.Category = challenge.Category
	if challenge.Difficulty.Valid {
		chall.Difficulty = challenge.Difficulty.String
	}
	chall.Instance = challenge.Type != sqlc.DeployTypeNormal
	chall.Hidden = challenge.Hidden
	chall.Points = challenge.Points
	chall.Solves = challenge.Solves
	chall.Tags = []string{}
	if tags != nil {
		chall.Tags = tags
	}
	chall.Solved = solved

	return &chall, nil
}

func GetChallenges(ctx context.Context, uid int32, author bool) ([]*Chall, error) {
	var challIDs []int32
	var err error
	challIDs, err = db.Sql.GetChallenges(ctx)
	if err != nil {
		return nil, err
	}

	challs := make([]*Chall, 0)
	for _, id := range challIDs {
		chall, err := GetChallengePreview(ctx, id, uid, author)
		if err != nil {
			return nil, err
		}
		if chall == nil {
			continue
		}
		challs = append(challs, chall)
	}

	return challs, nil
}
