package challenges_create

import (
	"context"
	"trxd/db"
	"trxd/db/sqlc"

	"github.com/lib/pq"
)

func CreateChallenge(ctx context.Context, name, category, description string,
	challType sqlc.DeployType, maxPoints int32, scoreType sqlc.ScoreType) (*sqlc.Challenge, error) {
	id, err := db.Sql.CreateChallenge(ctx, sqlc.CreateChallengeParams{
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

	return &sqlc.Challenge{
		ID:          id,
		Name:        name,
		Category:    category,
		Description: description,
		Type:        challType,
		MaxPoints:   maxPoints,
		ScoreType:   scoreType,
	}, nil
}
