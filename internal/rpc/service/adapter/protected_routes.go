package adapter

import (
	"os"

	"github.com/labstack/echo/v4"
	mw "github.com/loloneme/potential-waffle/internal/middleware"
)

func RegProtectedAdminRoutes(e *echo.Echo) {
	admin := e.Group("")
	admin.Use(mw.AdminAuthenticationMiddleware(os.Getenv("API_KEY")))

	admin.POST("/pullRequest/create", nil)
	admin.POST("/pullRequest/merge", nil)
	admin.POST("/pullRequest/reassign", nil)
	admin.POST("/users/setIsActive", nil)
	admin.POST("/users/bulkDeactivate", nil)
}
