package httpserver

import (
	goJson "encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/goccy/go-json"
	"github.com/labstack/echo/v4"
)

// CustomSerializer implements the echo.JSONSerializer.

type NewGoccyJsonEncoderFunc func(w io.Writer) *json.Encoder
type NewGoJsonEncoderFunc func(w io.Writer) *goJson.Encoder

type EncoderFunc interface {
	Encode(w io.Writer) interface{}
}

func (f NewGoccyJsonEncoderFunc) Encode(w io.Writer) interface{} {
	return f(w)
}

func (f NewGoJsonEncoderFunc) Encode(w io.Writer) interface{} {
	return f(w)
}

type CustomSerializer struct {
	encoder EncoderFunc
}

func NewDefaultSerializer() *CustomSerializer {
	return &CustomSerializer{
		encoder: NewGoccyJsonEncoderFunc(func(w io.Writer) *json.Encoder {
			return json.NewEncoder(w)
		}),
	}
}

func NewBuiltInGoSerializer() *CustomSerializer {
	return &CustomSerializer{
		encoder: NewGoJsonEncoderFunc(func(w io.Writer) *goJson.Encoder {
			return goJson.NewEncoder(w)
		}),
	}
}

// Serialize converts an interface into a json and writes it to the response.
// You can optionally use the indent parameter to produce pretty JSONs.
func (s *CustomSerializer) Serialize(c echo.Context, i interface{}, indent string) error {
	encInterface := s.encoder.Encode(c.Response())

	switch enc := encInterface.(type) {
	case *json.Encoder:
		if indent != "" {
			enc.SetIndent("", indent)
		}
		return enc.Encode(i)
	case *goJson.Encoder:
		if indent != "" {
			enc.SetIndent("", indent)
		}
		return enc.Encode(i)
	default:
		return fmt.Errorf("unknown encoder type: %T", encInterface)
	}
}

// Deserialize reads a JSON from a request body and converts it into an interface.
func (s *CustomSerializer) Deserialize(c echo.Context, i interface{}) error {
	err := json.NewDecoder(c.Request().Body).Decode(i)
	var ute *json.UnmarshalTypeError
	var se *json.SyntaxError
	if errors.As(err, &ute) {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			fmt.Sprintf(
				"Unmarshal type error: expected=%v, got=%v, field=%v, offset=%v",
				ute.Type,
				ute.Value,
				ute.Field,
				ute.Offset,
			),
		).SetInternal(err)
	} else if errors.As(err, &se) {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			fmt.Sprintf("Syntax error: offset=%v, error=%v", se.Offset, se.Error()),
		).SetInternal(err)
	}
	return err
}
