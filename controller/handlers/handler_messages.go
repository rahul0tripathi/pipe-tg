package handlers

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) MakeFetchAllMessages(svc MessageLogger) echo.HandlerFunc {
	return func(c echo.Context) error {
		return h.wrapper.WithContext(c, func(ctx context.Context) error {
			msg, err := svc.All(ctx)
			if err != nil {
				return err
			}

			return responseJSON(c, http.StatusOK, msg)
		})
	}
}
