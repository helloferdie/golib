package libmiddleware

import (
	"github.com/helloferdie/golib/libecho"
	"github.com/helloferdie/golib/libresponse"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// VerifyHeaderSecret - Verify header secret configuration
type VerifyHeaderSecretConfig struct {
	Skipper middleware.Skipper
	Field   string
	Value   string
}

// VerifyHeaderSecret - Verify header secret
func VerifyHeaderSecret(config VerifyHeaderSecretConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if config.Skipper != nil && config.Skipper(c) {
				return next(c)
			}

			if c.Request().Header.Get(config.Field) == config.Value {
				return next(c)
			}

			response := libresponse.GetDefault().ErrorUnauthorized()
			response.Error = "common.error.header.secret.invalid"
			return libecho.ParseResponse(c, response)
		}
	}
}
