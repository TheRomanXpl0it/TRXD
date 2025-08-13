package challenges_get

import (
	"context"
	"trxd/api/routes/challenge_get"
	"trxd/db"
)

func GetChallenges(ctx context.Context, uid int32, author bool) ([]*challenge_get.Chall, error) {
	var challIDs []int32
	var err error
	challIDs, err = db.Sql.GetChallenges(ctx)
	if err != nil {
		return nil, err
	}

	challs := make([]*challenge_get.Chall, 0)
	for _, id := range challIDs {
		chall, err := challenge_get.GetChallenge(ctx, id, uid, author)
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
