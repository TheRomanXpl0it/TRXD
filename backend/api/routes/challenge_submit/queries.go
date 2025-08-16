package challenge_submit

import (
	"context"
	"trxd/db"
	"trxd/db/sqlc"
)

func SubmitFlag(ctx context.Context, userID int32, challengeID int32, flag string) (sqlc.SubmissionStatus, error) {
	valid, err := db.Sql.CheckFlags(ctx, sqlc.CheckFlagsParams{
		Flag:    flag,
		ChallID: challengeID,
	})
	if err != nil {
		return sqlc.SubmissionStatusInvalid, err
	}

	status := sqlc.SubmissionStatusWrong
	if valid {
		status = sqlc.SubmissionStatusCorrect
	}

	status, err = db.Sql.Submit(ctx, sqlc.SubmitParams{
		UserID:  userID,
		ChallID: challengeID,
		Status:  status,
		Flag:    flag,
	})
	if err != nil {
		return sqlc.SubmissionStatusInvalid, err
	}

	return status, nil
}
