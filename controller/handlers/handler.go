package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/rahul0tripathi/pipetg/internal/integrations/tg"
)

type Handler struct {
	wrapper *tg.Client
}

func New(wrapper *tg.Client) *Handler {
	return &Handler{wrapper: wrapper}
}

type DataResponse struct {
	Data any `json:"data"`
}

func responseJSON(c echo.Context, status int, response interface{}) error {
	return c.JSON(
		status, DataResponse{
			Data: response,
		},
	)
}
