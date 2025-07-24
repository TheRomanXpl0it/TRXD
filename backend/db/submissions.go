package db

import "context"

func SubmitFlag(ctx context.Context, userID int32, challengeID int32, flag string) (SubmissionStatus, error) {
	valid, err := queries.CheckFlags(ctx, CheckFlagsParams{
		Flag:    flag,
		ChallID: challengeID,
	})
	if err != nil {
		return SubmissionStatusInvalid, err
	}

	status := SubmissionStatusWrong
	if valid {
		status = SubmissionStatusCorrect
	}

	status, err = queries.Submit(ctx, SubmitParams{
		UserID:  userID,
		ChallID: challengeID,
		Status:  status,
		Flag:    flag,
	})
	if err != nil {
		return SubmissionStatusInvalid, err
	}

	return status, nil
}
