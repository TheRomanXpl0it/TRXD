package submissions_create

import (
	"context"
	"trxd/db"
	"trxd/db/sqlc"
)

func SubmitFlag(ctx context.Context, userID int32, role sqlc.UserRole,
	challengeID int32, flag string) (sqlc.SubmissionStatus, bool, error) {
	valid, err := db.Sql.CheckFlags(ctx, sqlc.CheckFlagsParams{
		Flag:    flag,
		ChallID: challengeID,
	})
	if err != nil {
		return sqlc.SubmissionStatusInvalid, false, err
	}

	status := sqlc.SubmissionStatusWrong
	if valid {
		status = sqlc.SubmissionStatusCorrect
	}

	if role != sqlc.UserRolePlayer {
		return status, false, nil
	}

	res, err := db.Sql.Submit(ctx, sqlc.SubmitParams{
		UserID: userID,
		ID:     challengeID,
		Status: status,
		Flag:   flag,
	})
	if err != nil {
		return sqlc.SubmissionStatusInvalid, false, err
	}

	return res.Status, res.FirstBlood, nil
}
