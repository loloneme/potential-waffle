package middleware

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	echomw "github.com/oapi-codegen/echo-middleware"
)

func NewOpenAPIMiddleware(specPath string) echo.MiddlewareFunc {
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile(specPath)
	if err != nil {
		fmt.Errorf("failed to load OpenAPI spec: %v", err)
	}
	return echomw.OapiRequestValidator(doc)
}
