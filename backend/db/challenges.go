package db

import (
	"context"
	"database/sql"
	"strings"

	"github.com/lib/pq"
)

func CreateChallenge(ctx context.Context, name, category, description string,
	challType DeployType, maxPoints int32, scoreType ScoreType) (*Challenge, error) {
	id, err := queries.CreateChallenge(ctx, CreateChallengeParams{
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

	return &Challenge{
		ID:          id,
		Name:        name,
		Category:    category,
		Description: description,
		Type:        challType,
		MaxPoints:   maxPoints,
		ScoreType:   scoreType,
	}, nil
}

func GetChallengeByID(ctx context.Context, challengeID int32) (*Challenge, error) {
	challenge, err := queries.GetChallengeByID(ctx, challengeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &challenge, nil
}

func CreateFlag(ctx context.Context, challengeID int32, flag string, regex bool) (*Flag, error) {
	err := queries.CreateFlag(ctx, CreateFlagParams{
		Flag:    flag,
		ChallID: challengeID,
		Regex:   regex,
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // Unique violation error code
				return nil, nil
			}
		}
		return nil, err
	}

	return &Flag{
		ChallID: challengeID,
		Flag:    flag,
		Regex:   regex,
	}, nil
}

func GetFlagsByChallenge(ctx context.Context, challengeID int32) ([]GetFlagsByChallengeRow, error) {
	flags, err := queries.GetFlagsByChallenge(ctx, challengeID)
	if err != nil {
		return nil, err
	}

	return flags, nil
}

func GetTagsByChallenge(ctx context.Context, challengeID int32) ([]string, error) {
	tags, err := queries.GetTagsByChallenge(ctx, challengeID)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

type Chall struct {
	ID          int32                    `json:"id"`
	Name        string                   `json:"name"`
	Category    string                   `json:"category"`
	Description string                   `json:"description"`
	Difficulty  string                   `json:"difficulty"`
	Authors     []string                 `json:"authors"`
	Instance    bool                     `json:"instance"`
	Hidden      bool                     `json:"hidden"`
	Points      int32                    `json:"points"`
	Solves      int32                    `json:"solves"`
	Host        string                   `json:"host"`
	Port        int32                    `json:"port"`
	Attachments []string                 `json:"attachments"`
	Tags        []string                 `json:"tags"`
	Flags       []GetFlagsByChallengeRow `json:"flags"`
	Timeout     int32                    `json:"timeout"`
	Solved      bool                     `json:"solved"`
}

func GetChallenge(ctx context.Context, id int32, uid int32, author bool) (*Chall, error) {
	var chall Chall

	challenge, err := GetChallengeByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if challenge == nil {
		return nil, nil
	}
	if !author && challenge.Hidden {
		return nil, nil
	}

	flags, err := GetFlagsByChallenge(ctx, challenge.ID)
	if err != nil {
		return nil, err
	}

	tags, err := GetTagsByChallenge(ctx, challenge.ID)
	if err != nil {
		return nil, err
	}

	solved, err := queries.IsChallengeSolved(ctx, IsChallengeSolvedParams{
		ChallID: id,
		ID:      uid,
	})
	if err != nil {
		return nil, err
	}

	chall.ID = challenge.ID
	chall.Name = challenge.Name
	chall.Category = challenge.Category
	chall.Description = challenge.Description
	if challenge.Difficulty.Valid {
		chall.Difficulty = challenge.Difficulty.String
	}
	chall.Authors = []string{}
	if challenge.Authors.Valid {
		chall.Authors = strings.Split(challenge.Authors.String, ",") // TODO: change char
	}
	chall.Instance = challenge.Type != DeployTypeNormal
	chall.Hidden = challenge.Hidden
	chall.Points = challenge.Points
	chall.Solves = challenge.Solves
	if challenge.Host.Valid {
		chall.Host = challenge.Host.String
	}
	if challenge.Port.Valid {
		chall.Port = challenge.Port.Int32
	}
	chall.Attachments = []string{}
	if challenge.Attachments.Valid {
		chall.Attachments = strings.Split(challenge.Attachments.String, ",") // TODO: change char
	}
	if author {
		if flags != nil {
			chall.Flags = flags
		} else {
			chall.Flags = []GetFlagsByChallengeRow{}
		}
	}
	chall.Tags = []string{}
	if tags != nil {
		chall.Tags = tags
	}
	chall.Solved = solved

	return &chall, nil
}

func GetChallenges(ctx context.Context, uid int32, author bool) ([]*Chall, error) {
	var challIDs []int32
	var err error
	challIDs, err = queries.GetChallenges(ctx)
	if err != nil {
		return nil, err
	}

	challs := make([]*Chall, 0)
	for _, id := range challIDs {
		chall, err := GetChallenge(ctx, id, uid, author)
		if err != nil {
			return nil, err
		}
		if chall == nil {
			continue
		}
		challs = append(challs, chall)
	}

	return challs, nil
}

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

func DeleteChallenge(ctx context.Context, challengeID int32) error {
	err := queries.DeleteChallenge(ctx, challengeID)
	if err != nil {
		return err
	}

	return nil
}

func DeleteFlag(ctx context.Context, challengeID int32, flag string) error {
	err := queries.DeleteFlag(ctx, DeleteFlagParams{
		ChallID: challengeID,
		Flag:    flag,
	})
	if err != nil {
		return err
	}

	return nil
}
