package challenges_all_get

import (
	"context"
	"time"
	"trxd/db"
	"trxd/db/sqlc"
)

type Chall struct {
	ID           int32          `json:"id"`
	Name         string         `json:"name"`
	Category     string         `json:"category"`
	Description  string         `json:"description"`
	Difficulty   string         `json:"difficulty"`
	Authors      []string       `json:"authors"`
	Instance     bool           `json:"instance"`
	Hidden       bool           `json:"hidden"`
	Points       int            `json:"points"`
	Solves       int            `json:"solves"`
	Solved       bool           `json:"solved"`
	FirstBlood   bool           `json:"first_blood"`
	Attachments  []string       `json:"attachments"`
	Tags         []string       `json:"tags"`
	Host         string         `json:"host"`
	Port         int            `json:"port"`
	ConnType     sqlc.ConnType  `json:"conn_type"`
	MaxPoints    int            `json:"max_points"`
	ScoreType    sqlc.ScoreType `json:"score_type"`
	Timeout      int            `json:"timeout"`
	InstanceHost string         `json:"instance_host,omitempty"`
	InstancePort int            `json:"instance_port,omitempty"`
}

func GetChallenges(ctx context.Context, uid int32, tid int32, author bool) ([]Chall, error) {
	challenges, err := db.Sql.GetAllChallengesInfo(ctx, uid)
	if err != nil {
		return nil, err
	}

	challsData := make([]Chall, 0)
	for _, challenge := range challenges {
		if !author && challenge.Hidden {
			continue
		}

		chall := Chall{
			ID:          challenge.ID,
			Name:        challenge.Name,
			Category:    challenge.Category,
			Description: challenge.Description,
			Difficulty:  challenge.Difficulty,
			Authors:     challenge.Authors,
			Instance:    challenge.Type != sqlc.DeployTypeNormal,
			Hidden:      challenge.Hidden,
			Points:      int(challenge.Points),
			Solves:      int(challenge.Solves),
			Solved:      challenge.Solved,
			FirstBlood:  challenge.FirstBlood,
			Attachments: []string{},
			Tags:        []string{},
			Host:        challenge.Host,
			Port:        int(challenge.Port),
			ConnType:    challenge.ConnType,
			MaxPoints:   int(challenge.MaxPoints),
			ScoreType:   challenge.ScoreType,
			Timeout:     0,
		}

		if challenge.Attachments != nil {
			chall.Attachments = challenge.Attachments
		}
		if challenge.Tags != nil {
			chall.Tags = challenge.Tags
		}

		if challenge.ExpiresAt.Valid {
			chall.Timeout = int(time.Until(challenge.ExpiresAt.Time).Seconds())
			if chall.Timeout < 0 {
				chall.Timeout = 0
			}
		}
		if challenge.DockerID.Valid {
			if challenge.InstanceHost.Valid {
				chall.InstanceHost = challenge.InstanceHost.String
			}
			if challenge.InstancePort.Valid {
				chall.InstancePort = int(challenge.InstancePort.Int32)
			}
		}

		challsData = append(challsData, chall)
	}

	return challsData, nil
}
