package db

import "context"

func GetChallengeByID(ctx context.Context, challengeID int32) (*Challenge, error) {
	challenge, err := queries.GetChallengeByID(ctx, challengeID)
	if err != nil {
		return nil, err
	}

	return &challenge, nil
}
