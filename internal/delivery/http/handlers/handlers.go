package handlers

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jaam8/wb_tech_school_l0/internal/service"
	errs "github.com/jaam8/wb_tech_school_l0/pkg/errors"
)

type Handler struct {
	service *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) GetOrderById(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).
			JSON(fiber.Map{"error": "missing order id"})
	}

	order, err := h.service.GetOrder(c.UserContext(), id)
	if err != nil {
		if errors.Is(err, errs.ErrOrderNotFound) {
			return c.Status(http.StatusNotFound).
				JSON(fiber.Map{"error": "order not found"})
		}

		return c.Status(http.StatusInternalServerError).
			JSON(fiber.Map{"error": "internal server error"})
	}

	return c.Status(http.StatusOK).JSON(order)
}
