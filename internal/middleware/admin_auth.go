package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func AdminAuthenticationMiddleware(expected string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("X-Admin-Token")
			if token == "" || token != expected {
				return c.JSON(http.StatusUnauthorized, "token is empty")
			}
			return next(c)
		}
	}
}
