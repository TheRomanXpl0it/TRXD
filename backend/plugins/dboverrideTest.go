package main

import (
	"context"
	"strings"

	"trxd/pluginsapi/core"
	"trxd/pluginsapi/postgres"
)


func RegisterQueries(r postgres.Registry) {
	r.OnQuery("GetAllChallengesInfo", core.HandlerFunc(hookGetAllChallengesInfo), 0)
}

func hookGetAllChallengesInfo(ctx context.Context, args any) (any, error) {
	ev, ok := args.(postgres.QueryEvent)
	if !ok {
		return args, nil
	}
	ev.SQL = strings.Replace(ev.SQL, "c.name", "c.name || ' [PLUGIN]'", 1)
	return ev, nil
}
