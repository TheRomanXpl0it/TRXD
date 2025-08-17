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

func GetChallenges(ctx context.Context, uid int32, author bool) ([]*Chall, error) {
	challPreviews, err := db.Sql.GetChallengesPreview(ctx, uid)
	if err != nil {
		return nil, err
	}

	challsData := make([]*Chall, 0)
	for _, challenge := range challPreviews {
		if !author && challenge.Hidden {
			continue
		}

		tags, err := db.GetTagsByChallenge(ctx, challenge.ID)
		if err != nil {
			return nil, err
		}

		chall := &Chall{
			ID:       challenge.ID,
			Name:     challenge.Name,
			Category: challenge.Category,
			Instance: challenge.Type != sqlc.DeployTypeNormal,
			Hidden:   challenge.Hidden,
			Points:   challenge.Points,
			Solves:   challenge.Solves,
			Tags:     []string{},
			Solved:   challenge.Solved,
		}
		if challenge.Difficulty.Valid {
			chall.Difficulty = challenge.Difficulty.String
		}
		if tags != nil {
			chall.Tags = tags
		}

		challsData = append(challsData, chall)
	}

	return challsData, nil
}
