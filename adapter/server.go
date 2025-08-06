package adapter

import (
	"bytes"
	"encoding/json"

	"github.com/gofiber/fiber/v3"
)

func Server() ServerManager {
	return &manager{fiber.New()}
}

type ServerManager interface {
	Post(endpoint string, handler Handler)
	Get(endpoint string, handler Handler)
	Listen(addr string) error
}

type RequestContext interface {
	Query(name string, defaultValue ...string) string
	BodyParser(v any) error
	JSON(data any) error
	Status(statusCode int)RequestContext
	Send([]byte) error
	SendString(string) error
	SendStatus(statusCode int) error
}

type Handler func(RequestContext) error

func (h Handler) Fiber() fiber.Handler {
	return func(c fiber.Ctx) error {
		return h(reqContextAdapter{c})
	}
}

func (m *manager) Listen(addr string) error {
	return m.app.Listen(addr)
}

func (m *manager) Post(endpoint string, handler Handler) {
	m.app.Post(endpoint, handler.Fiber())
}

func (m *manager) Get(endpoint string, handler Handler) {
	m.app.Get(endpoint, handler.Fiber())
}

type manager struct {
	app *fiber.App
}

type reqContextAdapter struct {
	fiber.Ctx
}

func (c reqContextAdapter) BodyParser(v any) error {
	return json.NewDecoder(bytes.NewBuffer(c.Body())).Decode(v)
}

func (c reqContextAdapter) JSON(data any) error {
	return c.Ctx.JSON(data)
}

func (c reqContextAdapter) Status(statusCode int) RequestContext {
	return reqContextAdapter{c.Ctx.Status(statusCode)}
}