package submissions_create

import (
	"context"
	"fmt"
	"trxd/db"
	"trxd/db/sqlc"
)

func SubmitFlag(ctx context.Context, userID int32, challengeID int32, flag string) (sqlc.SubmissionStatus, bool, error) {
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

	res, err := db.Sql.Submit(ctx, sqlc.SubmitParams{
		UserID: userID,
		ID:     challengeID,
		Status: status,
		Flag:   flag,
	})
	if err != nil {
		return sqlc.SubmissionStatusInvalid, false, err
	}

	if !res.FirstBlood.Valid {
		return sqlc.SubmissionStatusInvalid, false, fmt.Errorf("unexpected null value for first_blood: %v", res.FirstBlood)
	}

	return res.Status, res.FirstBlood.Bool, nil
}
