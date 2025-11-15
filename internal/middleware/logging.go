package middleware

import (
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
)

func SlogMiddleware(logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			start := time.Now()

			err := next(c)

			logger.Info("http_request",
				slog.String("method", req.Method),
				slog.String("path", req.URL.Path),
				slog.Int("status", res.Status),
				slog.Duration("latency", time.Since(start)),
				slog.String("ip", c.RealIP()),
			)

			return err
		}
	}
}
