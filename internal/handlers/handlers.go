package handlers

import (
	"github.com/gofiber/fiber/v2"
)

type service interface {
	Get(id string) ([]byte, error)
}

type Handlers struct {
	service
}

func New(service service) *Handlers {
	return &Handlers{service}
}

func (h *Handlers) Get(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	result, err := h.service.Get(id)
	if err != nil {
		//TODO
	}
	return ctx.Send(result)
}
