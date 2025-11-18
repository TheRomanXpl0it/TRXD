package postgres

import (
	"context"

	"trxd/pluginsapi/core"
)

type QueryEvent struct {
	Name string
	SQL  string
	Args []any
}

type DBManager struct {
	queryHandlers map[string][]core.Handler
    triggerHandlers map[string][]core.Handler
}

type Registry interface {
    OnQuery(hook string, handler core.Handler, priority int)
    OnTrigger(hook string, handler core.Handler, priority int)
}

func NewDBManager() *DBManager {
    return &DBManager{
        queryHandlers: make(map[string][]core.Handler),
        triggerHandlers: make(map[string][]core.Handler),
    }
}

func (dm *DBManager) OnQuery(hook string, handler core.Handler, priority int) {
 	if dm.queryHandlers == nil {
        dm.queryHandlers = make(map[string][]core.Handler)
    }
	list := dm.queryHandlers[hook]
	list = core.InsertHandlerByPriority(list, handler, priority)
	dm.queryHandlers[hook] = list
}

func (dm *DBManager) OnTrigger(hook string, handler core.Handler, priority int) {
 	if dm.triggerHandlers == nil {
        dm.triggerHandlers = make(map[string][]core.Handler)
    }
	list := dm.triggerHandlers[hook]
	list = core.InsertHandlerByPriority(list, handler, priority)
	dm.triggerHandlers[hook] = list
}

func (dm *DBManager) DispatchQuery(ctx context.Context, ev QueryEvent) (QueryEvent, error) {
	args := ev
	if list, ok := dm.queryHandlers[ev.Name]; ok {
		for _, handler := range list {
			out, err := handler.Handle(ctx, args)
			if err != nil {
				return ev, err
			}
			if updated, ok := out.(QueryEvent); ok {
				ev = updated
				args = ev
			}
		}
	}
	
	return ev, nil
}