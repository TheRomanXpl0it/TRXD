package pluginsApi

import (
	"context"
	"fmt"
	"path/filepath"
	"plugin"
	
    "github.com/tde-nico/log"
	"trxd/pluginsapi/events"
	"trxd/pluginsapi/postgres"
	"trxd/pluginsapi/routes"
)

var Pm *Manager

type Manager struct {
    Events  	*events.EventManager
    Postgres 	*postgres.DBManager
    Routes      *routes.RoutesManager
}

func InitPlugins() error {
	eventsM := events.NewEventManager()
	postgresM := postgres.NewDBManager()
	routesM := routes.NewRoutesRoutesManager()
	
	Pm = &Manager {
		Events: eventsM,
		Postgres: postgresM,
		Routes: routesM,
	}
	
	matches, err := filepath.Glob("plugins/*.so")
	if err != nil {
		return err
	}

	for _, path := range matches {
		log.Infof("Loading plugin: %s", path)

		p, err := plugin.Open(path)
		if err != nil {
			log.Errorf("  -> failed to open: %v", err)
			continue
		}

		if sym, err := p.Lookup("RegisterEvents"); err == nil {
			if regFn, ok := sym.(func(events.Registry)); ok {
				regFn(Pm.Events)
				log.Infof("Loaded plugin %s event handler",path)
			} else {
				log.Errorf("  -> RegisterEvents has wrong type in %s", path)
			}
		}

		if sym, err := p.Lookup("RegisterDBs"); err == nil {
			if regFn, ok := sym.(func(postgres.Registry)); ok {
				regFn(Pm.Postgres)
				log.Infof("Loaded plugin %s DB handler",sym)
			} else {
				log.Errorf("  -> RegisterDB has wrong type in %s", path)
			}
		}

		if sym, err := p.Lookup("RegisterRoutes"); err == nil {
			if regFn, ok := sym.(func(routes.Registry)); ok {
				regFn(Pm.Routes)
				log.Infof("Loaded plugin %s Route handler",sym)
			} else {
				log.Errorf("  -> RegisterRoutes has wrong type in %s", path)
			}
		}
	}

	return nil
    
}

func  DispatchEvent[T any](ctx context.Context, hook string, args T) (T, error) {
	var zero T 
    out, err := Pm.Events.Dispatch(ctx, hook, args)
	if err != nil {
		return zero, err
	}
	res, ok := out.(T)
	if !ok {
		return zero, fmt.Errorf("plugin chain for %s returned incompatible type %T", hook, out)
	}
	return res, nil
}

func DispatchDb[T any](ctx context.Context, hook string, args T) (T, error) {
	var zero T 
    out, err := Pm.Postgres.Dispatch(ctx, hook, args)
	if err != nil {
		return zero, err
	}
	res, ok := out.(T)
	if !ok {
		return zero, fmt.Errorf("plugin chain for %s returned incompatible type %T", hook, out)
	}
	return res, nil
}