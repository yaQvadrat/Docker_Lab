package httpapi

import (
	"errors"

	"github.com/labstack/echo/v4"
)

var (
	ErrInternalServer    = errors.New("internal server error")
	ErrInvalidParameters = errors.New("invalid request parameters")
)

func newErrReasonJSON(c echo.Context, code int, msg interface{}) error {
	return c.JSON(code, map[string]interface{}{"reason": msg})
}
