package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) MakeFetchAllMessages(svc MessageLogger) echo.HandlerFunc {
	return func(c echo.Context) error {

		msg, err := svc.All(c.Request().Context())
		if err != nil {
			return err
		}

		return responseJSON(c, http.StatusOK, msg)
	}
}
