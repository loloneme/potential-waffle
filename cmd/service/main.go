package service

import (
	"context"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/loloneme/potential-waffle/internal"
	generated "github.com/loloneme/potential-waffle/internal/generated/openapi"
	"github.com/loloneme/potential-waffle/internal/rpc/service/adapter"
)

func main() {
	ctx := context.Background()

	db, err := internal.NewDatabaseConnection(ctx)
	if err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}
	defer db.Close()

	serviceAdapter := adapter.NewAdapter()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.User(validation.OpenAPIValidation())

	generated.RegisterHandlers(e, serviceAdapter)
}
