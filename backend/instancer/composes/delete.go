package composes

import (
	"context"

	"github.com/docker/compose/v2/pkg/api"
)

func KillCompose(ctx context.Context, name string) error {
	if ComposeCli == nil {
		return nil
	}

	err := ComposeCli.Down(ctx, name, api.DownOptions{})
	if err != nil {
		return err
	}

	return nil
}
