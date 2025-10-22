package teams_scoreboard_graph

import (
	"context"
	"sort"
	"time"
	"trxd/db"
)

type Submission struct {
	ChallID    int       `json:"chall_id"`
	Score      int       `json:"score"`
	FirstBlood bool      `json:"first_blood"`
	Timestamp  time.Time `json:"timestamp"`
}

type Top struct {
	TeamID      int          `json:"team_id"`
	Submissions []Submission `json:"submissions"`
}

func QueryTeamScoreboardGraph(ctx context.Context) ([]Top, error) {
	res, err := db.Sql.GetTeamsScoreboardGraph(ctx)
	if err != nil {
		return nil, err
	}

	top := make(map[int32]Top, 0)
	for _, row := range res {
		if _, ok := top[row.TeamID]; !ok {
			top[row.TeamID] = Top{
				TeamID:      int(row.TeamID),
				Submissions: make([]Submission, 0),
			}
		}

		data := top[row.TeamID]
		prev := 0
		if len(data.Submissions) > 0 {
			prev = data.Submissions[len(data.Submissions)-1].Score
		}

		data.Submissions = append(data.Submissions, Submission{
			ChallID:    int(row.ChallID),
			Score:      prev + int(row.Points),
			FirstBlood: row.FirstBlood,
			Timestamp:  row.Timestamp,
		})

		top[row.TeamID] = data
	}

	topArr := make([]Top, 0, len(top))
	for _, v := range top {
		topArr = append(topArr, v)
	}
	sort.Slice(topArr, func(i, j int) bool {
		subsA, subsB := topArr[i].Submissions, topArr[j].Submissions
		if len(subsA) == 0 || len(subsB) == 0 {
			return len(subsA) > len(subsB)
		}
		a, b := subsA[len(subsA)-1], subsB[len(subsB)-1]
		if a.Score == b.Score {
			return a.Timestamp.Before(b.Timestamp)
		}
		return a.Score > b.Score
	})

	return topArr, nil
}
