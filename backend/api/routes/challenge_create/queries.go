package challenge_create

import (
	"context"
	"trxd/db"

	"github.com/lib/pq"
)

func CreateChallenge(ctx context.Context, name, category, description string,
	challType db.DeployType, maxPoints int32, scoreType db.ScoreType) (*db.Challenge, error) {
	id, err := db.Sql.CreateChallenge(ctx, db.CreateChallengeParams{
		Name:        name,
		Category:    category,
		Description: description,
		Type:        challType,
		MaxPoints:   maxPoints,
		ScoreType:   scoreType,
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // Unique violation error code
				return nil, nil
			}
		}
		return nil, err
	}

	return &db.Challenge{
		ID:          id,
		Name:        name,
		Category:    category,
		Description: description,
		Type:        challType,
		MaxPoints:   maxPoints,
		ScoreType:   scoreType,
	}, nil
}
