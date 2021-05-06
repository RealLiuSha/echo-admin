package echox

import (
	"net/http"

	"github.com/RealLiuSha/echo-admin/errors"

	"github.com/labstack/echo/v4"
)

// Response in order to unify the returned response structure
type Response struct {
	Code    int         `json:"-"`
	Pretty  bool        `json:"-"`
	Data    interface{} `json:"data,omitempty"`
	Message interface{} `json:"message"`
}

// sends a JSON response with status code.
func (a Response) JSON(ctx echo.Context) error {
	if a.Message == "" || a.Message == nil {
		a.Message = http.StatusText(a.Code)
	}

	if err, ok := a.Message.(error); ok {
		if errors.Is(err, errors.DatabaseInternalError) {
			a.Code = http.StatusInternalServerError
		}

		if errors.Is(err, errors.DatabaseRecordNotFound) {
			a.Code = http.StatusNotFound
		}

		a.Message = err.Error()
	}

	if a.Pretty {
		return ctx.JSONPretty(a.Code, a, "\t")
	}

	return ctx.JSON(a.Code, a)
}
