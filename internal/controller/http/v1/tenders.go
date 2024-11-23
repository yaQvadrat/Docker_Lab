package httpapi

import (
	"app/internal/service"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type tenderRoutes struct {
	tenderService service.Tender
}

func newTenderRoutes(s service.Tender) *tenderRoutes {
	return &tenderRoutes{s}
}

type NewTenderDTO struct {
	Name            string    `json:"name" validate:"required,max=100"`
	Description     string    `json:"description" validate:"required,max=500"`
	ServiceType     string    `json:"serviceType" validate:"required,oneof=Construction Delivery Manufacture"`
	OrganizationId  uuid.UUID `json:"organizationId" validate:"required"`
	CreatorUsername string    `json:"creatorUsername" validate:"required,max=50"`
}

func (r *tenderRoutes) newTender(c echo.Context) error {
	// Binding and validation
	var input NewTenderDTO
	if err := c.Bind(&input); err != nil {
		return handleBindingError(c, err)
	}

	if err := c.Validate(input); err != nil {
		return newErrReasonJSON(c, http.StatusBadRequest, err.Error())
	}

	// Creating tender
	tender, err := r.tenderService.CreateTender(c.Request().Context(), service.CreateTenderInput{
		Name:            input.Name,
		Description:     input.Description,
		ServiceType:     input.ServiceType,
		OrganizationId:  input.OrganizationId,
		CreatorUsername: input.CreatorUsername,
	})
	if err != nil {
		if errors.Is(err, service.ErrUsername) {
			return newErrReasonJSON(c, http.StatusUnauthorized, err.Error())
		}
		if errors.Is(err, service.ErrForbidden) {
			return newErrReasonJSON(c, http.StatusForbidden, err.Error())
		}
		return newErrReasonJSON(c, http.StatusInternalServerError, err.Error())
	}

	// Create response
	type response struct {
		Id          uuid.UUID `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Status      string    `json:"status"`
		ServiceType string    `json:"serviceType"`
		Version     int       `json:"version"`
		CreatedAt   string    `json:"createdAt"`
	}

	return c.JSON(http.StatusOK, response{
		Id:          tender.Id,
		Name:        tender.Name,
		Description: tender.Description,
		Status:      tender.Status,
		ServiceType: tender.Type,
		Version:     tender.Version,
		CreatedAt:   tender.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

type MyTendersDTO struct {
	LimitAndOffset
	Username string `query:"username" validate:"required,max=50"`
}

func (r *tenderRoutes) myTenders(c echo.Context) error {
	// Binding and validation
	var input MyTendersDTO
	if err := c.Bind(&input); err != nil {
		return handleBindingError(c, err)
	}

	if err := c.Validate(input); err != nil {
		return newErrReasonJSON(c, http.StatusBadRequest, err.Error())
	}
	if err := limitAndOffsetValidate(&input.LimitAndOffset); err != nil {
		return newErrReasonJSON(c, http.StatusBadRequest, err.Error())
	}

	// Get tender list by username
	tenders, err := r.tenderService.GetTendersByUsername(c.Request().Context(), service.GetByUsernameInput{
		Limit:    int(input.Limit.Int32),
		Offset:   int(input.Offset.Int32),
		Username: input.Username,
	})
	if err != nil {
		if errors.Is(err, service.ErrUsername) {
			return newErrReasonJSON(c, http.StatusUnauthorized, err.Error())
		}
		return newErrReasonJSON(c, http.StatusInternalServerError, err.Error())
	}

	// Create response
	type response struct {
		Id          uuid.UUID `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Status      string    `json:"status"`
		ServiceType string    `json:"serviceType"`
		Version     int       `json:"version"`
		CreatedAt   string    `json:"createdAt"`
	}
	responseBatch := []response{}
	for _, t := range tenders {
		responseBatch = append(responseBatch, response{
			Id:          t.Id,
			Name:        t.Name,
			Description: t.Description,
			Status:      t.Status,
			ServiceType: t.Type,
			Version:     t.Version,
			CreatedAt:   t.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return c.JSON(http.StatusOK, responseBatch)
}

type TendersDTO struct {
	LimitAndOffset
	ServiceType []string `query:"service_type" validate:"dive,oneof=Construction Delivery Manufacture"`
}

func (r *tenderRoutes) tenders(c echo.Context) error {
	// Binding and validation
	var input TendersDTO
	if err := c.Bind(&input); err != nil {
		return handleBindingError(c, err)
	}

	if err := c.Validate(input); err != nil {
		return newErrReasonJSON(c, http.StatusBadRequest, err.Error())
	}
	if err := limitAndOffsetValidate(&input.LimitAndOffset); err != nil {
		return newErrReasonJSON(c, http.StatusBadRequest, err.Error())
	}

	// Get tender list
	tenders, err := r.tenderService.GetTenders(c.Request().Context(), service.GetTendersInput{
		Limit:       int(input.Limit.Int32),
		Offset:      int(input.Offset.Int32),
		ServiceType: input.ServiceType,
	})
	if err != nil {
		return newErrReasonJSON(c, http.StatusInternalServerError, err.Error())
	}

	// Create response
	type response struct {
		Id          uuid.UUID `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Status      string    `json:"status"`
		ServiceType string    `json:"serviceType"`
		Version     int       `json:"version"`
		CreatedAt   string    `json:"createdAt"`
	}
	responseBatch := []response{}
	for _, t := range tenders {
		responseBatch = append(responseBatch, response{
			Id:          t.Id,
			Name:        t.Name,
			Description: t.Description,
			Status:      t.Status,
			ServiceType: t.Type,
			Version:     t.Version,
			CreatedAt:   t.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return c.JSON(http.StatusOK, responseBatch)
}

type PutStatusDTO struct {
	TenderId uuid.UUID `param:"tenderId" validate:"required"`
	Status   string    `query:"status" validate:"required,oneof=Created Published Closed"`
	Username string    `query:"username" validate:"required,max=50"`
}

func (r *tenderRoutes) putStatus(c echo.Context) error {
	// Binding and validation
	var input PutStatusDTO
	if err := c.Bind(&input); err != nil {
		return handleBindingError(c, err)
	}
	queryBinder := &echo.DefaultBinder{}
	if err := queryBinder.BindQueryParams(c, &input); err != nil {
		return handleBindingError(c, err)
	}

	if err := c.Validate(input); err != nil {
		return newErrReasonJSON(c, http.StatusBadRequest, err.Error())
	}

	// Change status
	tender, err := r.tenderService.ChangeStatus(c.Request().Context(), service.ChangeTenderStatusInput{
		TenderId: input.TenderId,
		Status:   input.Status,
		Username: input.Username,
	})
	if err != nil {
		if errors.Is(err, service.ErrUsername) {
			return newErrReasonJSON(c, http.StatusUnauthorized, err.Error())
		}
		if errors.Is(err, service.ErrNotFoundTender) {
			return newErrReasonJSON(c, http.StatusNotFound, err.Error())
		}
		if errors.Is(err, service.ErrForbidden) {
			return newErrReasonJSON(c, http.StatusForbidden, err.Error())
		}
		return newErrReasonJSON(c, http.StatusInternalServerError, err.Error())
	}

	// Create response
	type response struct {
		Id          uuid.UUID `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Status      string    `json:"status"`
		ServiceType string    `json:"serviceType"`
		Version     int       `json:"version"`
		CreatedAt   string    `json:"createdAt"`
	}

	return c.JSON(http.StatusOK, response{
		Id:          tender.Id,
		Name:        tender.Name,
		Description: tender.Description,
		Status:      tender.Status,
		ServiceType: tender.Type,
		Version:     tender.Version,
		CreatedAt:   tender.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

type GetStatusDTO struct {
	TenderId uuid.UUID `param:"tenderId" validate:"required"`
	Username string    `query:"username" validate:"required,max=50"`
}

func (r *tenderRoutes) getStatus(c echo.Context) error {
	// Binding and validation
	var input GetStatusDTO
	if err := c.Bind(&input); err != nil {
		return handleBindingError(c, err)
	}

	if err := c.Validate(input); err != nil {
		return newErrReasonJSON(c, http.StatusBadRequest, err.Error())
	}

	// Get status
	tender, err := r.tenderService.GetTender(c.Request().Context(), input.TenderId, input.Username)
	if err != nil {
		if errors.Is(err, service.ErrUsername) {
			return newErrReasonJSON(c, http.StatusUnauthorized, err.Error())
		}
		if errors.Is(err, service.ErrNotFoundTender) {
			return newErrReasonJSON(c, http.StatusNotFound, err.Error())
		}
		if errors.Is(err, service.ErrForbidden) {
			return newErrReasonJSON(c, http.StatusForbidden, err.Error())
		}
		return newErrReasonJSON(c, http.StatusInternalServerError, err.Error())
	}

	// Create response
	return c.JSON(http.StatusOK, tender.Status)
}

type EditTenderDTO struct {
	TenderId uuid.UUID `param:"tenderId" validate:"required"`
	Username string    `query:"username" validate:"required,max=50"`
	EditTenderBody
}

func (r *tenderRoutes) editTender(c echo.Context) error {
	// Binding and validation
	var input EditTenderDTO
	if err := c.Bind(&input); err != nil {
		return handleBindingError(c, err)
	}
	queryBinder := &echo.DefaultBinder{}
	if err := queryBinder.BindQueryParams(c, &input); err != nil {
		return handleBindingError(c, err)
	}

	if err := c.Validate(input); err != nil {
		return newErrReasonJSON(c, http.StatusBadRequest, err.Error())
	}
	if err := editTenderBodyValidate(&input.EditTenderBody); err != nil {
		return newErrReasonJSON(c, http.StatusBadRequest, err.Error())
	}

	// Edit tender
	tender, err := r.tenderService.Edit(c.Request().Context(), service.EditTenderInput{
		TenderId:    input.TenderId,
		Username:    input.Username,
		Name:        input.Name.String,
		Description: input.Description.String,
		ServiceType: input.ServiceType.String,
	})
	if err != nil {
		if errors.Is(err, service.ErrUsername) {
			return newErrReasonJSON(c, http.StatusUnauthorized, err.Error())
		}
		if errors.Is(err, service.ErrNotFoundTender) {
			return newErrReasonJSON(c, http.StatusNotFound, err.Error())
		}
		if errors.Is(err, service.ErrForbidden) {
			return newErrReasonJSON(c, http.StatusForbidden, err.Error())
		}
		return newErrReasonJSON(c, http.StatusInternalServerError, err.Error())
	}

	// Create response
	type response struct {
		Id          uuid.UUID `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Status      string    `json:"status"`
		ServiceType string    `json:"serviceType"`
		Version     int       `json:"version"`
		CreatedAt   string    `json:"createdAt"`
	}

	return c.JSON(http.StatusOK, response{
		Id:          tender.Id,
		Name:        tender.Name,
		Description: tender.Description,
		Status:      tender.Status,
		ServiceType: tender.Type,
		Version:     tender.Version,
		CreatedAt:   tender.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

type RollbackTenderDTO struct {
	TenderId uuid.UUID `param:"tenderId" validate:"required"`
	Version  int       `param:"version" validate:"required,gt=0"`
	Username string    `query:"username" validate:"required,max=50"`
}

func (r *tenderRoutes) rollbackTender(c echo.Context) error {
	// Binding and validation
	var input RollbackTenderDTO
	if err := c.Bind(&input); err != nil {
		return handleBindingError(c, err)
	}
	queryBinder := &echo.DefaultBinder{}
	if err := queryBinder.BindQueryParams(c, &input); err != nil {
		return handleBindingError(c, err)
	}

	if err := c.Validate(input); err != nil {
		return newErrReasonJSON(c, http.StatusBadRequest, err.Error())
	}

	// Rollback tender
	tender, err := r.tenderService.Rollback(c.Request().Context(), service.RollbackTenderInput{
		TenderId: input.TenderId,
		Version:  input.Version,
		Username: input.Username,
	})
	if err != nil {
		if errors.Is(err, service.ErrUsername) {
			return newErrReasonJSON(c, http.StatusUnauthorized, err.Error())
		}
		if errors.Is(err, service.ErrNotFoundTender) {
			return newErrReasonJSON(c, http.StatusNotFound, err.Error())
		}
		if errors.Is(err, service.ErrForbidden) {
			return newErrReasonJSON(c, http.StatusForbidden, err.Error())
		}
		return newErrReasonJSON(c, http.StatusInternalServerError, err.Error())
	}

	// Create response
	type response struct {
		Id          uuid.UUID `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Status      string    `json:"status"`
		ServiceType string    `json:"serviceType"`
		Version     int       `json:"version"`
		CreatedAt   string    `json:"createdAt"`
	}

	return c.JSON(http.StatusOK, response{
		Id:          tender.Id,
		Name:        tender.Name,
		Description: tender.Description,
		Status:      tender.Status,
		ServiceType: tender.Type,
		Version:     tender.Version,
		CreatedAt:   tender.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}
