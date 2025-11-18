package events

import (
    "context"
    "trxd/pluginsapi/core"

    "github.com/tde-nico/log"
)

type EventManager struct {
    handlers map[string][]core.Handler
}

type Registry interface {
    OnEvent(hook string, handler core.Handler, priority int)
}

func NewEventManager() *EventManager {
    return &EventManager{
        handlers: make(map[string][]core.Handler),
    }
}

func (em *EventManager) OnEvent(hook string, handler core.Handler, priority int) {
    if em.handlers == nil {
        em.handlers = make(map[string][]core.Handler)
    }
    list := em.handlers[hook]
    list = core.InsertHandlerByPriority(list, handler, priority)
    em.handlers[hook] = list
}

func (em *EventManager) Dispatch(ctx context.Context, hook string, args any) (any, error) {
    handlers, ok := em.handlers[hook]
    if !ok || len(handlers) == 0 {
        return args, nil
    }

    in := args
    var err error

    for _, handler := range handlers {
        beforeCall := in

        in, err = handler.Handle(ctx, in)
        if typeErr := core.CheckSameTypes(beforeCall, in); typeErr != nil {
            log.Errorf("plugin for %s returned incompatible types in safe mode: %v", hook, typeErr)
            in = beforeCall
            continue
        }
        

        if err != nil {
            log.Error("Error in executing plugin", "hook", hook, "err", err)
            in = beforeCall
        }
    }

    return in, nil
}
