package httpserver

import (
	"fmt"
	"runtime"

	"github.com/labstack/echo/v4"

	"github.com/labstack/echo/v4/middleware"
)

// RecoverOption is a function type that accepts a Recover pointer object
// and update it.
type RecoverOption func(rec *Recover)

// Recover is a mechanism to recover from panics and print the
// panic error in logs.
type Recover struct {
	// Skipper defines a function to skip middleware.
	Skipper middleware.Skipper

	// Size of the stack to be printed.
	// Optional. Default value 4KB.
	StackSize int

	// DisableStackAll disables formatting stack traces of all other goroutines
	// into buffer after the trace for the current goroutine.
	// Optional. Default value false.
	DisableStackAll bool

	// DisablePrintStack disables printing stack trace.
	// Optional. Default value as false.
	DisablePrintStack bool
}

func newDefaultRecover() *Recover {
	return &Recover{
		Skipper:           middleware.DefaultSkipper,
		StackSize:         4 << 10, // 4 KB
		DisableStackAll:   false,
		DisablePrintStack: false,
	}
}

func WithSkipper(skipper middleware.Skipper) RecoverOption {
	return func(config *Recover) {
		config.Skipper = skipper
	}
}

func WithStackSize(stackSize int) RecoverOption {
	return func(config *Recover) {
		config.StackSize = stackSize
	}
}

func WithDisableStackAll(disable bool) RecoverOption {
	return func(config *Recover) {
		config.DisableStackAll = disable
	}
}

func WithDisablePrintStack(disable bool) RecoverOption {
	return func(config *Recover) {
		config.DisableStackAll = disable
	}
}

func RecoverMiddleware(options ...RecoverOption) echo.MiddlewareFunc {
	// Defaults
	rec := newDefaultRecover()
	for _, option := range options {
		option(rec)
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if rec.Skipper(c) {
				return next(c)
			}

			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					stack := make([]byte, rec.StackSize)
					length := runtime.Stack(stack, !rec.DisableStackAll)
					if !rec.DisablePrintStack {
						c.Logger().Printf("[PANIC RECOVER] %v %s\n", err, stack[:length])
					}

					c.Error(err)
				}
			}()
			return next(c)
		}
	}
}
