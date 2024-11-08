package handlers

import "github.com/labstack/echo/v4"

type Handler struct{}

func New() *Handler {
	return &Handler{}
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
