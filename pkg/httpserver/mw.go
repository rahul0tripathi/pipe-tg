package httpserver

import (
	"bytes"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func MakeLoggerMiddleware() echo.MiddlewareFunc {
	return middleware.LoggerWithConfig(
		middleware.LoggerConfig{
			CustomTagFunc: func(c echo.Context, buf *bytes.Buffer) (int, error) {
				return buf.WriteString(fmt.Sprintf(`"URL":"%s"`, c.Request().URL.String()))
			},
			Format: `{ "level":"info", "request_id":"${request_id}","component":"echo", "host":"${host}${url}", "sourceAddress":"${ip}:${port}", "method":"${method}", "status":"${status}", "path": "${path}", "trackingId": "${reqHeader:TrackingId}" }` + "\n",
		},
	)
}

func MakeCorsMiddleware() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(
		middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{"GET", "POST", "DELETE", "PATCH"},
		},
	)
}
