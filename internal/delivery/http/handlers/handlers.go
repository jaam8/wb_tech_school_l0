package handlers

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jaam8/wb_tech_school_l0/internal/delivery/http/schemas"
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

// GetOrderById godoc
// @Summary get order by id
// @Description returns an order by its order_uid
// @Tags order
// @Accept json
// @Produce json
// @Success 200 {object} models.Order
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 404 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /api/v1/order/{id} [get]
func (h *Handler) GetOrderById(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).
			JSON(schemas.ErrorResponse{Error: "missing order id"})
	}

	order, err := h.service.GetOrder(c.UserContext(), id)
	if err != nil {
		if errors.Is(err, errs.ErrOrderNotFound) {
			return c.Status(http.StatusNotFound).
				JSON(schemas.ErrorResponse{Error: err.Error()})
		}

		return c.Status(http.StatusInternalServerError).
			JSON(schemas.ErrorResponse{Error: errs.ErrInternalServerError.Error()})
	}

	return c.Status(http.StatusOK).JSON(order)
}

// Ping godoc
// @Summary health checker
// @Description Returns "pong" if the service is alive
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {string} string "pong"
// @Router /ping [get]
func Ping(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON("pong")
}
