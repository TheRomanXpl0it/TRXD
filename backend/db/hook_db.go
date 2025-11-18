package db

import (
	"context"
	"database/sql"
	"strings"
	
	"trxd/pluginsapi"
	"trxd/pluginsapi/postgres"
	
	"github.com/tde-nico/log"
)

type HookedDB struct {
	Inner *sql.DB
}

func extractSqlcName(query string) string {
	lines := strings.SplitN(query, "\n", 2)
	if len(lines) == 0 {
		return ""
	}
	line := strings.TrimSpace(lines[0])
	if !strings.HasPrefix(line, "-- name:") {
		return ""
	}

	rest := strings.TrimSpace(strings.TrimPrefix(line, "-- name:"))
	parts := strings.Fields(rest)
	if len(parts) == 0 {
		return ""
	}
	return parts[0]
}

func (h *HookedDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	ev := postgres.QueryEvent{
		Name: extractSqlcName(query),
		SQL:  query,
		Args: args,
	}

	if pluginsApi.Pm != nil && pluginsApi.Pm.Postgres != nil {
		var err error
		ev, err = pluginsApi.Pm.Postgres.DispatchQuery(ctx, ev)
		if err != nil {
			return nil, err
		}
	}

	return h.Inner.ExecContext(ctx, ev.SQL, ev.Args...)
}

func (h *HookedDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	ev := postgres.QueryEvent{
		Name: extractSqlcName(query),
		SQL:  query,
		Args: args,
	}

	if pluginsApi.Pm != nil && pluginsApi.Pm.Postgres != nil {
		var err error
		ev, err = pluginsApi.Pm.Postgres.DispatchQuery(ctx, ev)
		if err != nil {
			return nil, err
		}
	}

	return h.Inner.QueryContext(ctx, ev.SQL, ev.Args...)
}

func (h *HookedDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	ev := postgres.QueryEvent{
		Name: extractSqlcName(query),
		SQL:  query,
		Args: args,
	}

	if pluginsApi.Pm != nil && pluginsApi.Pm.Postgres != nil {
		var err error
		ev, err = pluginsApi.Pm.Postgres.DispatchQuery(ctx, ev)
		if err != nil {
			log.Error("Plugin error during query:","err",err,"query",query)
		}
	}

	return h.Inner.QueryRowContext(ctx, ev.SQL, ev.Args...)
}

func (h *HookedDB) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return h.Inner.PrepareContext(ctx, query)
}