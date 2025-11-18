package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Registry interface {
	Handle(method, path string, h fiber.Handler) error
	Replace(method, path string, h fiber.Handler) error
}

type RoutesManager struct {
	routes map[string]fiber.Handler
}

func NewRoutesRoutesManager() *RoutesManager {
	return &RoutesManager{
		routes: make(map[string]fiber.Handler),
	}
}

func key(method, path string) string {
	return method + " " + path
}

func (m *RoutesManager) Handle(method, path string, h fiber.Handler) error {
	k := key(method, path)
	if _, exists := m.routes[k]; exists {
		return fmt.Errorf("Handler for the path %s already exists.",path)
	}
	m.routes[k] = h
	return nil
}

// Replace overwrites any existing route only if it already exists.
func (m *RoutesManager) Replace(method, path string, h fiber.Handler) error {
	k := key(method, path)
	if _, exists := m.routes[k]; exists {
		m.routes[k] = h
		return nil
	}
	return fmt.Errorf("Handler for the path %s does not exist.",path)
}

func (m *RoutesManager) Mount(app *fiber.App) {
	app.All("/*", func(c *fiber.Ctx) error {
		k := key(c.Method(), c.Path())

		if handler, ok := m.routes[k]; ok {
			return handler(c)
		}

		return c.Next()
	})
}
