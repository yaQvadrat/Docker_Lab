package httpapi

import (
	"errors"
	"net/http"

	"github.com/guregu/null/v5"
	"github.com/labstack/echo/v4"
)

const (
	DefaultLimit  = 5
	DefaultOffset = 0
)

type LimitAndOffset struct {
	Limit  null.Int32 `query:"limit"`
	Offset null.Int32 `query:"offset"`
}

func limitAndOffsetValidate(input *LimitAndOffset) error {
	if input.Limit.Valid && input.Limit.Int32 < 0 {
		return ErrInvalidParameters
	}
	if input.Offset.Valid && input.Offset.Int32 < 0 {
		return ErrInvalidParameters
	}

	if !input.Limit.Valid {
		input.Limit.Int32 = DefaultLimit
	}
	if !input.Offset.Valid {
		input.Offset.Int32 = DefaultOffset
	}
	return nil
}

type EditTenderBody struct {
	Name        null.String `json:"name"`
	Description null.String `json:"description"`
	ServiceType null.String `json:"serviceType"`
}

func editTenderBodyValidate(input *EditTenderBody) error {
	if !input.Name.Valid && !input.Description.Valid && !input.ServiceType.Valid {
		return ErrInvalidParameters
	}
	if input.Name.Valid && !input.Description.Valid && !input.ServiceType.Valid {
		if len(input.Name.String) > 100 || len(input.Name.String) <= 0 {
			return ErrInvalidParameters
		}
		return nil
	}
	if input.Name.Valid && input.Description.Valid && !input.ServiceType.Valid {
		if len(input.Name.String) > 100 || len(input.Name.String) <= 0 {
			return ErrInvalidParameters
		}
		if len(input.Description.String) > 500 || len(input.Description.String) <= 0 {
			return ErrInvalidParameters
		}
		return nil
	}
	if input.Name.Valid && input.Description.Valid && input.ServiceType.Valid {
		if len(input.Name.String) > 100 || len(input.Name.String) <= 0 {
			return ErrInvalidParameters
		}
		if len(input.Description.String) > 500 || len(input.Description.String) <= 0 {
			return ErrInvalidParameters
		}
		if input.ServiceType.String != "Construction" && input.ServiceType.String != "Delivery" && input.ServiceType.String != "Manufacture" {
			return ErrInvalidParameters
		}
		return nil
	}
	if input.Name.Valid && !input.Description.Valid && input.ServiceType.Valid {
		if len(input.Name.String) > 100 || len(input.Name.String) <= 0 {
			return ErrInvalidParameters
		}
		if input.ServiceType.String != "Construction" && input.ServiceType.String != "Delivery" && input.ServiceType.String != "Manufacture" {
			return ErrInvalidParameters
		}
		return nil
	}
	if !input.Name.Valid && input.Description.Valid && input.ServiceType.Valid {
		if len(input.Description.String) > 500 || len(input.Description.String) <= 0 {
			return ErrInvalidParameters
		}
		if input.ServiceType.String != "Construction" && input.ServiceType.String != "Delivery" && input.ServiceType.String != "Manufacture" {
			return ErrInvalidParameters
		}
		return nil
	}
	if !input.Name.Valid && input.Description.Valid && !input.ServiceType.Valid {
		if len(input.Description.String) > 500 || len(input.Description.String) <= 0 {
			return ErrInvalidParameters
		}
		return nil
	}
	if !input.Name.Valid && !input.Description.Valid && input.ServiceType.Valid {
		if input.ServiceType.String != "Construction" && input.ServiceType.String != "Delivery" && input.ServiceType.String != "Manufacture" {
			return ErrInvalidParameters
		}
		return nil
	}
	return ErrInvalidParameters
}

type EditBidBody struct {
	Name        null.String `json:"name"`
	Description null.String `json:"description"`	
}

func editBidBodyValidate(input *EditBidBody) error {
	if !input.Name.Valid && !input.Description.Valid {
		return ErrInvalidParameters
	}
	if input.Name.Valid && input.Description.Valid {
		if len(input.Name.String) > 100 || len(input.Name.String) <= 0 {
			return ErrInvalidParameters
		}
		if len(input.Description.String) > 500 || len(input.Description.String) <= 0 {
			return ErrInvalidParameters
		}
		return nil
	}
	if input.Name.Valid && !input.Description.Valid {
		if len(input.Name.String) > 100 || len(input.Name.String) <= 0 {
			return ErrInvalidParameters
		}
		return nil
	}
	if !input.Name.Valid && input.Description.Valid {
		if len(input.Description.String) > 500 || len(input.Description.String) <= 0 {
			return ErrInvalidParameters
		}
		return nil
	}
	return ErrInvalidParameters
}

func handleBindingError(c echo.Context, err error) error {
	var httpErr *echo.HTTPError
	if ok := errors.As(err, &httpErr); ok {
		return newErrReasonJSON(c, http.StatusBadRequest, httpErr.Message)
	}
	return newErrReasonJSON(c, http.StatusBadRequest, ErrInvalidParameters.Error())
}
