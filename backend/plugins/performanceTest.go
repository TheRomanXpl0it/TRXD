package main

import (
	"context"

	"trxd/api/routes/challenges_all_get"
	"trxd/pluginsapi/core"
	"trxd/pluginsapi/events"
)

func RegisterEvents(reg events.Registry) {
	reg.OnEvent("challengesGet", core.HandlerFunc(onChallengesGet), 1)
}

func onChallengesGet(ctx context.Context, payload any) (any, error) {
	challs, ok := payload.([]challenges_all_get.Chall)
	if !ok {
		print("Plugin error!")
		return payload, nil
	}

	for i := range challs {
		if challs[i].Category == "Web" {
			challs[i].Name = challs[i].Name + " [WEB PLUGIN]"
		}
	}

	return challs, nil
}