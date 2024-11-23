package httpapi

import (
	"app/internal/service"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type bidRoutes struct {
	bidService service.Bid
}

func newBidRoutes(s service.Bid) *bidRoutes {
	return &bidRoutes{s}
}

type NewBidDTO struct {
	Name        string    `json:"name" validate:"required,max=100"`
	Description string    `json:"description" validate:"required,max=500"`
	TenderId    uuid.UUID `json:"tenderId" validate:"required"`
	AuthorType  string    `json:"authorType" validate:"required,oneof=User Organization"`
	AuthorId    uuid.UUID `json:"authorId" validate:"required"`
}

func (r *bidRoutes) newBid(c echo.Context) error {
	// Binding and validation
	var input NewBidDTO
	if err := c.Bind(&input); err != nil {
		return handleBindingError(c, err)
	}

	if err := c.Validate(input); err != nil {
		return newErrReasonJSON(c, http.StatusBadRequest, err.Error())
	}

	// Create bid
	bid, err := r.bidService.CreateBid(c.Request().Context(), service.CreateBidInput{
		Name:        input.Name,
		Description: input.Description,
		TenderId:    input.TenderId,
		AuthorType:  input.AuthorType,
		AuthorId:    input.AuthorId,
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
		Id         uuid.UUID `json:"id"`
		Name       string    `json:"name"`
		Status     string    `json:"status"`
		AuthorType string    `json:"authorType"`
		AuthorId   uuid.UUID `json:"authorId"`
		Version    int       `json:"version"`
		CreatedAt  string    `json:"createdAt"`
	}

	return c.JSON(http.StatusOK, response{
		Id:         bid.Id,
		Name:       bid.Name,
		Status:     bid.Status,
		AuthorType: bid.AuthorType,
		AuthorId:   bid.AuthorId,
		Version:    bid.Version,
		CreatedAt:  bid.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

type SubmitDecisionDTO struct {
	BidId    uuid.UUID `param:"bidId" validate:"required"`
	Decision string    `query:"decision" validate:"required,oneof=Approved Rejected"`
	Username string    `query:"username" validate:"required,max=50"`
}

func (r *bidRoutes) submitDecision(c echo.Context) error {
	// Binding and validation
	var input SubmitDecisionDTO
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

	// Make decision
	bid, err := r.bidService.SubmitDecision(c.Request().Context(), input.BidId, input.Username, input.Decision)
	if err != nil {
		if errors.Is(err, service.ErrUsername) {
			return newErrReasonJSON(c, http.StatusUnauthorized, err.Error())
		}
		if errors.Is(err, service.ErrNotFoundTender) {
			return newErrReasonJSON(c, http.StatusBadRequest, err.Error())
		}
		if errors.Is(err, service.ErrNotFoundBid) {
			return newErrReasonJSON(c, http.StatusNotFound, err.Error())
		}
		if errors.Is(err, service.ErrForbidden) {
			return newErrReasonJSON(c, http.StatusForbidden, err.Error())
		}
		return newErrReasonJSON(c, http.StatusInternalServerError, err.Error())
	}

	// Create response
	type response struct {
		Id         uuid.UUID `json:"id"`
		Name       string    `json:"name"`
		Status     string    `json:"status"`
		AuthorType string    `json:"authorType"`
		AuthorId   uuid.UUID `json:"authorId"`
		Version    int       `json:"version"`
		CreatedAt  string    `json:"createdAt"`
	}

	return c.JSON(http.StatusOK, response{
		Id:         bid.Id,
		Name:       bid.Name,
		Status:     bid.Status,
		AuthorType: bid.AuthorType,
		AuthorId:   bid.AuthorId,
		Version:    bid.Version,
		CreatedAt:  bid.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

type PutBidStatusDTO struct {
	BidId    uuid.UUID `param:"bidId" validate:"required"`
	Status   string    `query:"status" validate:"required,oneof=Created Published Canceled"`
	Username string    `query:"username" validate:"required,max=50"`
}

func (r *bidRoutes) putStatus(c echo.Context) error {
	// Binding and validation
	var input PutBidStatusDTO
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
	bid, err := r.bidService.ChangeStatus(c.Request().Context(), input.BidId, input.Status, input.Username)
	if err != nil {
		if errors.Is(err, service.ErrUsername) {
			return newErrReasonJSON(c, http.StatusUnauthorized, err.Error())
		}
		if errors.Is(err, service.ErrNotFoundBid) {
			return newErrReasonJSON(c, http.StatusNotFound, err.Error())
		}
		if errors.Is(err, service.ErrForbidden) {
			return newErrReasonJSON(c, http.StatusForbidden, err.Error())
		}
		return newErrReasonJSON(c, http.StatusInternalServerError, err.Error())
	}

	// Create response
	type response struct {
		Id         uuid.UUID `json:"id"`
		Name       string    `json:"name"`
		Status     string    `json:"status"`
		AuthorType string    `json:"authorType"`
		AuthorId   uuid.UUID `json:"authorId"`
		Version    int       `json:"version"`
		CreatedAt  string    `json:"createdAt"`
	}

	return c.JSON(http.StatusOK, response{
		Id:         bid.Id,
		Name:       bid.Name,
		Status:     bid.Status,
		AuthorType: bid.AuthorType,
		AuthorId:   bid.AuthorId,
		Version:    bid.Version,
		CreatedAt:  bid.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

type GetBidStatusDTO struct {
	BidId    uuid.UUID `param:"bidId" validate:"required"`
	Username string    `query:"username" validate:"required,max=50"`
}

func (r *bidRoutes) getStatus(c echo.Context) error {
	// Binding and validation
	var input GetBidStatusDTO
	if err := c.Bind(&input); err != nil {
		return handleBindingError(c, err)
	}

	if err := c.Validate(input); err != nil {
		return newErrReasonJSON(c, http.StatusBadRequest, err.Error())
	}

	// Get status
	bid, err := r.bidService.Get(c.Request().Context(), input.BidId, input.Username)
	if err != nil {
		if errors.Is(err, service.ErrUsername) {
			return newErrReasonJSON(c, http.StatusUnauthorized, err.Error())
		}
		if errors.Is(err, service.ErrNotFoundBid) {
			return newErrReasonJSON(c, http.StatusNotFound, err.Error())
		}
		if errors.Is(err, service.ErrForbidden) {
			return newErrReasonJSON(c, http.StatusForbidden, err.Error())
		}
		return newErrReasonJSON(c, http.StatusInternalServerError, err.Error())
	}

	// Create respone
	return c.JSON(http.StatusOK, bid.Status)
}

type EditBidDTO struct {
	BidId    uuid.UUID `param:"bidId" validate:"required"`
	Username string    `query:"username" validate:"required,max=50"`
	EditBidBody
}

func (r *bidRoutes) editBid(c echo.Context) error {
	// Binding and validation
	var input EditBidDTO
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
	if err := editBidBodyValidate(&input.EditBidBody); err != nil {
		return newErrReasonJSON(c, http.StatusBadRequest, err.Error())
	}

	// Edit bid
	bid, err := r.bidService.Edit(c.Request().Context(), service.EditBidInput{
		BidId:       input.BidId,
		Username:    input.Username,
		Name:        input.Name.String,
		Description: input.Description.String,
	})
	if err != nil {
		if errors.Is(err, service.ErrUsername) {
			return newErrReasonJSON(c, http.StatusUnauthorized, err.Error())
		}
		if errors.Is(err, service.ErrNotFoundBid) {
			return newErrReasonJSON(c, http.StatusNotFound, err.Error())
		}
		if errors.Is(err, service.ErrForbidden) {
			return newErrReasonJSON(c, http.StatusForbidden, err.Error())
		}
		return newErrReasonJSON(c, http.StatusInternalServerError, err.Error())
	}

	// Create response
	type response struct {
		Id         uuid.UUID `json:"id"`
		Name       string    `json:"name"`
		Status     string    `json:"status"`
		AuthorType string    `json:"authorType"`
		AuthorId   uuid.UUID `json:"authorId"`
		Version    int       `json:"version"`
		CreatedAt  string    `json:"createdAt"`
	}

	return c.JSON(http.StatusOK, response{
		Id:         bid.Id,
		Name:       bid.Name,
		Status:     bid.Status,
		AuthorType: bid.AuthorType,
		AuthorId:   bid.AuthorId,
		Version:    bid.Version,
		CreatedAt:  bid.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

type RollbackBidDTO struct {
	BidId    uuid.UUID `param:"bidId" validate:"required"`
	Version  int       `param:"version" validate:"required,gt=0"`
	Username string    `query:"username" validate:"required,max=50"`
}

func (r *bidRoutes) rollbackBid(c echo.Context) error {
	// Binding and validation
	var input RollbackBidDTO
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

	// Rollback bid
	bid, err := r.bidService.Rollback(c.Request().Context(), input.BidId, input.Version, input.Username)
	if err != nil {
		if errors.Is(err, service.ErrUsername) {
			return newErrReasonJSON(c, http.StatusUnauthorized, err.Error())
		}
		if errors.Is(err, service.ErrNotFoundBid) {
			return newErrReasonJSON(c, http.StatusNotFound, err.Error())
		}
		if errors.Is(err, service.ErrForbidden) {
			return newErrReasonJSON(c, http.StatusForbidden, err.Error())
		}
		return newErrReasonJSON(c, http.StatusInternalServerError, err.Error())
	}

	// Create response
	type response struct {
		Id         uuid.UUID `json:"id"`
		Name       string    `json:"name"`
		Status     string    `json:"status"`
		AuthorType string    `json:"authorType"`
		AuthorId   uuid.UUID `json:"authorId"`
		Version    int       `json:"version"`
		CreatedAt  string    `json:"createdAt"`
	}

	return c.JSON(http.StatusOK, response{
		Id:         bid.Id,
		Name:       bid.Name,
		Status:     bid.Status,
		AuthorType: bid.AuthorType,
		AuthorId:   bid.AuthorId,
		Version:    bid.Version,
		CreatedAt:  bid.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}
