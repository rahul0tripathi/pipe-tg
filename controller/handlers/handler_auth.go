package handlers

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) MakeHandleSendCode(svc AuthFlowSvc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return h.wrapper.WithUnAuthEchoContext(c, func(ctx context.Context) error {
			err := svc.RequestNewCode(ctx)
			if err != nil {
				return err
			}

			return responseJSON(c, http.StatusOK, "processing")
		})
	}
}

type submitCodeRequest struct {
	Code string `json:"code"`
}

func (h *Handler) MakeHandleSubmitCode(svc AuthFlowSvc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return h.wrapper.WithUnAuthEchoContext(c, func(ctx context.Context) error {
			req := &submitCodeRequest{}
			if err := c.Bind(req); err != nil {
				return err
			}

			err := svc.SubmitCode(ctx, req.Code)
			if err != nil {
				return err
			}

			return responseJSON(c, http.StatusOK, "authenticated")
		})
	}
}
