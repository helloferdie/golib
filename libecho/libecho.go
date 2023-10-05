package libecho

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/helloferdie/golib/liblogger"
	"github.com/helloferdie/golib/libresponse"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Initialize - Initialize required middleware before assign route
func Initialize(e *echo.Echo) {
	e.HTTPErrorHandler = ErrorHandler
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(Logger)
}

// StartHTTP - Start as HTTP server
func StartHTTP(e *echo.Echo) {
	// Start server
	go func() {
		if err := e.Start(":" + os.Getenv("app_port")); err != nil {
			liblogger.Log(nil, false).Error(err)
			liblogger.Log(nil, false).Error("Fail start HTTP server")
			liblogger.Log(nil, false).Error("Shutting down the server")
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		liblogger.Log(nil, false).Error("Fail shutting down server")
		os.Exit(1)
	} else {
		liblogger.Log(nil, false).Info("Shutdown HTTP server - done")
	}
}

// Logger - Log every incoming request
func Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log := map[string]interface{}{
			"method":   c.Request().Method,
			"uri":      c.Request().URL.String(),
			"ip":       c.Request().RemoteAddr,
			"ip_real":  c.Request().Header.Get("X-Real-Ip"),
			"ip_proxy": c.Request().Header.Get("X-Proxy-Ip"),
		}
		liblogger.Log(log, false).Info("incoming request")
		return next(c)
	}
}

// ErrorHandler - Custom error handler
func ErrorHandler(err error, c echo.Context) {
	response := libresponse.GetDefault()
	report, ok := err.(*echo.HTTPError)
	if !ok {
		report = echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response.Code = report.Code
	switch response.Code {
	case 400:
		response.Message = "common.error.request.default"
		reportMsg := report.Message.(string)
		if reportMsg == "missing or malformed jwt" {
			response.Error = "common.error.request.jwt.required"
		} else {
			if e1, ok := report.Internal.(*json.UnmarshalTypeError); ok {
				response.Error = "common.error.request.unmarshal_var"
				response.ErrorVar = map[string]interface{}{
					"f": e1.Field,
					"e": e1.Type,
					"v": e1.Value,
				}
			}
		}

		if response.Error == "" {
			response.Error = "common.error.request.bad"
		}
	case 401:
		response.Message = "common.error.request.default"
		reportMsg := report.Message.(string)
		if reportMsg == "invalid or expired jwt" {
			reportInternal := report.Internal.Error()
			if reportInternal == "Token is expired" {
				response.Error = "common.error.request.jwt.expired"
			} else {
				response.Error = "common.error.request.jwt.invalid"
			}
		}

		if response.Error == "" {
			response.Error = "common.error.request.unauthorized"
		}
	case 403:
		response.Message = "common.error.request.default"
		response.Error = "common.error.request.forbidden"
	case 404:
		response.Message = "common.error.request.default"
		if response.Error == "" {
			response.Error = "common.error.request.route.not_found"
		}
	case 405:
		response.Message = "common.error.request.default"
		response.Error = "common.error.request.method"
	case 413:
		response.Message = "common.error.request.default"
		response.Error = "common.error.request.size.large"
	case 415:
		response.Message = "common.error.request.default"
		response.Error = "common.error.request.media_type"
	case 422:
		response.Message = "validation.error.default"
	case 429:
		response.Message = "common.error.request.default"
		response.Error = "common.error.request.limit"
	case 500:
		response.Message = "common.error.server.internal"
		if response.Error == "" {
			response.Error = "common.error.server.internal"
		}
	case 502:
		response.Message = "common.error.server.gateway"
		response.Error = "common.error.service.unreachable"
	}

	ParseResponse(c, response)
}

// ParseResponse - Parse return response in JSON format
func ParseResponse(c echo.Context, response *libresponse.Response) (err error) {
	return c.JSON(response.Code, response)
}
