package httpserver

import (
	"context"

	"github.com/labstack/echo/v4"
)

type Router interface {
	Use(middleware ...echo.MiddlewareFunc)
	Pre(middleware ...echo.MiddlewareFunc)
	Group(prefix string, middleware ...echo.MiddlewareFunc) (sg *echo.Group)
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

type Server struct {
	app     *echo.Echo
	notify  chan error
	address string
}

func NewServer(address string) *Server {
	app := echo.New()

	app.JSONSerializer = NewDefaultSerializer()

	return &Server{
		app:     app,
		notify:  make(chan error),
		address: address,
	}
}

func NewServerWithMiddlewares(address string) *Server {
	app := echo.New()

	app.JSONSerializer = NewDefaultSerializer()

	app.Use(
		MakeCorsMiddleware(),
		MakeLoggerMiddleware(),
		RecoverMiddleware(),
	)

	return &Server{
		app:     app,
		notify:  make(chan error),
		address: address,
	}
}

func (s *Server) Use(middlewares ...echo.MiddlewareFunc) {
	for i := range middlewares {
		s.app.Use(middlewares[i])
	}
}

func (s *Server) Pre(middlewares ...echo.MiddlewareFunc) {
	for i := range middlewares {
		s.app.Pre(middlewares[i])
	}
}

func (s *Server) Router() Router {
	return s.app
}

func (s *Server) JSONSerializer(serializer echo.JSONSerializer) {
	s.app.JSONSerializer = serializer
}

func (s *Server) Start() {
	go func() {
		s.notify <- s.app.Start(s.address)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	return s.app.Shutdown(context.Background())
}

func ResponseJSON(c echo.Context, status int, response interface{}) error {
	return c.JSON(status, response)
}
