package db

import "context"

func SubmitFlag(userID int32, challengeID int32, flag string) (SubmissionStatus, error) {
	valid, err := queries.CheckFlags(context.Background(), CheckFlagsParams{
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

	status, err = queries.Submit(context.Background(), SubmitParams{
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
