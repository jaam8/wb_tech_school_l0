package middlewares

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jaam8/wb_tech_school_l0/pkg/logger"
	"go.uber.org/zap"
)

func LogMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, err := logger.New(c.UserContext())
		if err != nil {
			ctx, _ = logger.New(context.Background())
		}

		c.SetUserContext(ctx)

		logger.Info(ctx, "Request",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
		)
		err = c.Next()

		logger.Info(ctx, "Response",
			zap.Int("status", c.Response().StatusCode()),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
		)

		return err
	}
}
