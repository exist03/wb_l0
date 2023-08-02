package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"net/http"
	"wb_l0/common"
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
		if errors.Is(err, common.ErrNotFound) {
			return ctx.SendStatus(http.StatusNotFound)
		} else if errors.Is(err, common.ErrInvalidID) {
			return ctx.SendStatus(http.StatusBadRequest)
		} else {
			log.Err(err)
			return ctx.SendStatus(http.StatusInternalServerError)
		}
	}
	return ctx.Send(result)
}
