package teams_get

import (
	"context"
	"database/sql"
	"time"
	"trxd/db"
	"trxd/db/sqlc"
	"trxd/utils"
)

type Solve struct {
	ID         int32     `json:"id"`
	Name       string    `json:"name"`
	Category   string    `json:"category"`
	Points     int32     `json:"points"`
	FirstBlood bool      `json:"first_blood"`
	Timestamp  time.Time `json:"timestamp"`
	UserID     int32     `json:"user_id,omitempty"`
}

type TeamData struct {
	ID      int32                       `json:"id"`
	Name    string                      `json:"name"`
	Score   int32                       `json:"score"`
	Country string                      `json:"country"`
	Members []sqlc.GetTeamMembersRow    `json:"members,omitempty"`
	Solves  []Solve                     `json:"solves"`
	Badges  []sqlc.GetBadgesFromTeamRow `json:"badges"`
}

func getMembers(ctx context.Context, teamID int32, admin bool) ([]sqlc.GetTeamMembersRow, error) {
	allMembers, err := db.Sql.GetTeamMembers(ctx, sql.NullInt32{Int32: teamID, Valid: true})
	if err != nil {
		return nil, err
	}

	members := make([]sqlc.GetTeamMembersRow, 0)
	for _, member := range allMembers {
		if !admin && utils.In(member.Role, []sqlc.UserRole{sqlc.UserRoleAuthor, sqlc.UserRoleAdmin}) {
			continue
		}

		members = append(members, member)
	}

	return members, nil
}

func getSolves(ctx context.Context, teamID int32, userMode bool) ([]Solve, error) {
	solvesRaw, err := db.Sql.GetTeamSolves(ctx, teamID)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	solves := make([]Solve, 0, len(solvesRaw))
	for _, solveRaw := range solvesRaw {
		solve := Solve{
			ID:         solveRaw.ID,
			Name:       solveRaw.Name,
			Category:   solveRaw.Category,
			Points:     solveRaw.Points,
			FirstBlood: solveRaw.FirstBlood,
			Timestamp:  solveRaw.Timestamp,
		}

		if !userMode {
			solve.UserID = solveRaw.UserID
		}

		solves = append(solves, solve)
	}

	return solves, nil
}

func GetTeam(ctx context.Context, teamID int32, admin bool) (*TeamData, error) {
	teamData := TeamData{}

	team, err := db.GetTeamByID(ctx, teamID)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, nil
	}

	modeStr, err := db.GetConfig(ctx, "user-mode")
	if err != nil {
		return nil, err
	}
	userMode := modeStr == "true"

	teamData.ID = team.ID
	teamData.Name = team.Name
	teamData.Score = team.Score
	if team.Country.Valid {
		teamData.Country = team.Country.String
	}

	if !userMode {
		members, err := getMembers(ctx, teamID, admin)
		if err != nil {
			return nil, err
		}
		teamData.Members = members
	}

	solves, err := getSolves(ctx, teamID, userMode)
	if err != nil {
		return nil, err
	}
	teamData.Solves = solves

	badges, err := db.Sql.GetBadgesFromTeam(ctx, teamID)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}
	if badges == nil {
		badges = make([]sqlc.GetBadgesFromTeamRow, 0)
	}
	teamData.Badges = badges

	return &teamData, nil
}
