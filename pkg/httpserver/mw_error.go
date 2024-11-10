package httpserver

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rahul0tripathi/pipetg/pkg/log"
)

type Error struct {
	Error string `json:"errors"`
}

type ErrorResponse struct {
	Data Error `json:"data"`
}

func NewErrorResponse(err string) ErrorResponse {
	return ErrorResponse{
		Data: Error{Error: err},
	}
}

func makeErrorCatcherMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)

			// return if there's no error
			if err == nil {
				return nil
			}

			logger := log.GetLogger(c.Request().Context())
			logger.Error("error middleware", log.Err(err), log.Str("path", c.Path()))

			return c.JSON(http.StatusInternalServerError, NewErrorResponse(fmt.Sprintf("err: %s", err.Error())))

		}
	}
}
