package submissions_get

import (
	"context"
	"database/sql"
	"time"
	"trxd/db"
	"trxd/db/sqlc"
)

type GetSubmissionsUserMode struct {
	ID         int32                 `json:"id"`
	UserID     *int32                `json:"user_id,omitempty"`
	UserName   string                `json:"user_name,omitempty"`
	TeamID     int32                 `json:"team_id"`
	TeamName   string                `json:"team_name"`
	ChallID    int32                 `json:"chall_id"`
	ChallName  string                `json:"chall_name"`
	Status     sqlc.SubmissionStatus `json:"status"`
	FirstBlood bool                  `json:"first_blood"`
	Flag       string                `json:"flag"`
	Timestamp  time.Time             `json:"timestamp"`
}

func getSubs(ctx context.Context, offset int32, limit int32) (int64, []sqlc.GetSubmissionsRow, error) {
	total, err := db.Sql.GetTotalSubmissions(ctx)
	if err != nil {
		return 0, nil, err
	}

	submissions, err := db.Sql.GetSubmissions(ctx, sqlc.GetSubmissionsParams{
		Offset: offset,
		Limit:  sql.NullInt32{Int32: limit, Valid: limit != 0},
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return total, []sqlc.GetSubmissionsRow{}, nil
		}
		return 0, nil, err
	}

	return total, submissions, nil
}

func GetSubmissions(ctx context.Context, offset int32, limit int32) (int64, []GetSubmissionsUserMode, error) {
	userModeStr, err := db.GetConfig(ctx, "user-mode")
	if err != nil {
		return 0, nil, err
	}
	userMode := userModeStr == "true"

	total, submissions, err := getSubs(ctx, offset, limit)
	if err != nil {
		return 0, nil, err
	}

	subs := make([]GetSubmissionsUserMode, len(submissions))
	for i, submission := range submissions {
		subs[i] = GetSubmissionsUserMode{
			ID:         submission.ID,
			TeamID:     submission.TeamID,
			TeamName:   submission.TeamName,
			ChallID:    submission.ChallID,
			ChallName:  submission.ChallName,
			Status:     submission.Status,
			FirstBlood: submission.FirstBlood,
			Flag:       submission.Flag,
			Timestamp:  submission.Timestamp,
		}
		if !userMode {
			subs[i].UserID = &submission.UserID
			subs[i].UserName = submission.UserName
		}
	}

	return total, subs, nil
}
