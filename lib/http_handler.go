package lib

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/RealLiuSha/echo-admin/pkg/echox"
	"github.com/RealLiuSha/echo-admin/pkg/slice"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type HttpHandler struct {
	Engine   *echo.Echo
	RouterV1 *echo.Group

	Validate *validator.Validate
}

// HtppServer validation function
type Validator struct {
	validate *validator.Validate
}

// Implement the bind method to verify the request's struct for parameter validation
type BinderWithValidation struct{}

func (a *Validator) Validate(i interface{}) error {
	return a.validate.Struct(i)
}

// NewHttpHandler creates a new request handler
func NewHttpHandler(logger Logger, config Config) HttpHandler {
	// Error handlers
	echo.NotFoundHandler = func(ctx echo.Context) error {
		return echox.Response{Code: http.StatusNotFound}.JSON(ctx)
	}

	echo.MethodNotAllowedHandler = func(ctx echo.Context) error {
		return echox.Response{Code: http.StatusMethodNotAllowed}.JSON(ctx)
	}

	// new engine
	engine := echo.New()
	engine.HidePort = true
	engine.HideBanner = true
	engine.Binder = &BinderWithValidation{}

	// set http handler
	httpHandler := HttpHandler{
		Engine:   engine,
		RouterV1: engine.Group("/api/v1"),
	}

	// custom the error handler
	httpHandler.Engine.HTTPErrorHandler = func(err error, ctx echo.Context) {
		var (
			code    = http.StatusInternalServerError
			message interface{}
		)

		he, ok := err.(*echo.HTTPError)
		if ok {
			code = he.Code
			message = he.Message

			if he.Internal != nil {
				message = fmt.Errorf("%v - %v", message, he.Internal)
			}
		}

		// Send response
		if !ctx.Response().Committed {
			// https://www.w3.org/Protocols/rfc2616/rfc2616-sec9.html
			if ctx.Request().Method == http.MethodHead {
				err = ctx.NoContent(he.Code)
			} else {
				err = echox.Response{
					Code:    code,
					Message: message,
				}.JSON(ctx)
			}

			if err != nil {
				logger.DesugarZap.Error(err.Error())
			}
		}
	}

	// override the default validator
	httpHandler.Engine.Validator = func() echo.Validator {
		v := validator.New()

		v.RegisterValidation("json", func(fl validator.FieldLevel) bool {
			var js json.RawMessage
			return json.Unmarshal([]byte(fl.Field().String()), &js) == nil
		})

		v.RegisterValidation("in", func(fl validator.FieldLevel) bool {
			value := fl.Field().String()
			if slice.ContainsString(strings.Split(fl.Param(), ";"), value) || value == "" {
				return true
			}

			return false
		})

		return &Validator{validate: v}
	}()

	return httpHandler
}

func (BinderWithValidation) Bind(i interface{}, ctx echo.Context) error {
	binder := &echo.DefaultBinder{}

	if err := binder.Bind(i, ctx); err != nil {
		return errors.New(err.(*echo.HTTPError).Message.(string))
	}

	if err := ctx.Validate(i); err != nil {
		// Validate only provides verification function for struct.
		// When the requested data type is not struct,
		// the variable should be considered legal after the bind succeeds.
		if reflect.TypeOf(i).Kind() != reflect.Struct {
			return nil
		}

		var buf bytes.Buffer
		if ferrs, ok := err.(validator.ValidationErrors); ok {
			for _, ferr := range ferrs {
				buf.WriteString("Validation failed on ")
				buf.WriteString(ferr.Tag())
				buf.WriteString(" for ")
				buf.WriteString(ferr.StructField())
				buf.WriteString("\n")
			}

			return errors.New(buf.String())
		}

		return err
	}

	return nil
}
