package db

import "context"

func GetChallengeByID(challengeID int32) (*Challenge, error) {
	challenge, err := queries.GetChallengeByID(context.Background(), challengeID)
	if err != nil {
		return nil, err
	}

	return &challenge, nil
}
