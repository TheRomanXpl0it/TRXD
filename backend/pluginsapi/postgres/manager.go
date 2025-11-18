package postgres

import (
    "context"

    "github.com/tde-nico/log"
    "trxd/pluginsapi/core"
)

type DBManager struct {
    handlers map[string][]core.Handler
}

type Registry interface {
    OnDbHook(hook string, handler core.Handler, priority int)
}

func NewDBManager() *DBManager {
    return &DBManager{
        handlers: make(map[string][]core.Handler),
    }
}

func (dm *DBManager) OnDbHook(hook string, handler core.Handler, priority int) {
    if dm.handlers == nil {
        dm.handlers = make(map[string][]core.Handler)
    }
    list := dm.handlers[hook]
    list = core.InsertHandlerByPriority(list, handler, priority)
    dm.handlers[hook] = list
}

func (m *DBManager) Dispatch(ctx context.Context, hook string, args any) (any, error) {
    handlers, ok := m.handlers[hook]
    if !ok || len(handlers) == 0 {
        return args, nil
    }

    in := args
    var err error

    for _, handler := range handlers {
        beforeCall := in

        in, err = handler.Handle(ctx, in)
        if typeErr := core.CheckSameTypes(beforeCall, in); typeErr != nil {
            log.Errorf("db plugin for %s returned incompatible types in safe mode: %v", hook, typeErr)
            in = beforeCall
            continue
        }

        if err != nil {
            log.Error("Error in executing db plugin", "hook", hook, "err", err)
            in = beforeCall
        }
    }

    return in, nil
}