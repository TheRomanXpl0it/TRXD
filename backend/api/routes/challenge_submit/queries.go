package challenge_submit

import (
	"context"
	"trxd/db"
)

func SubmitFlag(ctx context.Context, userID int32, challengeID int32, flag string) (db.SubmissionStatus, error) {
	valid, err := db.Sql.CheckFlags(ctx, db.CheckFlagsParams{
		Flag:    flag,
		ChallID: challengeID,
	})
	if err != nil {
		return db.SubmissionStatusInvalid, err
	}

	status := db.SubmissionStatusWrong
	if valid {
		status = db.SubmissionStatusCorrect
	}

	status, err = db.Sql.Submit(ctx, db.SubmitParams{
		UserID:  userID,
		ChallID: challengeID,
		Status:  status,
		Flag:    flag,
	})
	if err != nil {
		return db.SubmissionStatusInvalid, err
	}

	return status, nil
}
