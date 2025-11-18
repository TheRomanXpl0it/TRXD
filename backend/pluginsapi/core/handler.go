package core

import "context"

type Handler interface {
    Handle(ctx context.Context, args any) (any, error)
}

type HandlerFunc func(ctx context.Context, args any) (any, error)

func (f HandlerFunc) Handle(ctx context.Context, args any) (any, error) {
    return f(ctx, args)
}

func InsertHandlerByPriority(handlers []Handler, h Handler, priority int) []Handler {

    if priority < 0 {
    priority = 0
    }
    if priority > len(handlers) {
        priority = len(handlers)
    }

    handlers = append(handlers, nil)
    copy(handlers[priority+1:], handlers[priority:])
    handlers[priority] = h

    return handlers
}