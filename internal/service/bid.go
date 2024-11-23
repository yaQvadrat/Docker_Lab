package service

import (
	e "app/internal/entity"
	"app/internal/repo"
	"app/internal/repo/repoerrors"
	rt "app/internal/repo/repotypes"
	"context"
	"errors"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type BidService struct {
	tenderRepo   repo.Tender
	employeeRepo repo.Employee
	bidRepo      repo.Bid
}

func NewBidService(tRepo repo.Tender, eRepo repo.Employee, bRepo repo.Bid) *BidService {
	return &BidService{
		tenderRepo:   tRepo,
		employeeRepo: eRepo,
		bidRepo:      bRepo,
	}
}

func (s *BidService) CreateBid(ctx context.Context, in CreateBidInput) (e.Bid, error) {
	// Check if user exists
	_, err := s.employeeRepo.GetById(ctx, in.AuthorId)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return e.Bid{}, ErrUsername
		}
		log.Errorf("BidService.Create - employeeRepo.GetByid: %v", err)
		return e.Bid{}, ErrGetEmployeeById
	}

	// Check if tender exists and published
	tender, err := s.tenderRepo.Get(ctx, in.TenderId, rt.VersionLatest)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return e.Bid{}, ErrNotFoundTender
		}
		log.Errorf("BidService.Create - tenderRepo.Get: %v", err)
		return e.Bid{}, ErrGetTender
	}
	if tender.Status != "Published" {
		return e.Bid{}, ErrNotFoundTender
	}

	// Check resposibility in case when AuthorType = "Organization"
	if in.AuthorType == "Organization" {
		isResponsible, err := s.employeeRepo.IsResponsibleSimplified(ctx, in.AuthorId)
		if err != nil {
			log.Errorf("BidService.Create - employeeRepo.IsResponsibleSimplified: %v", err)
			return e.Bid{}, ErrCheckResponsibility
		}
		if !isResponsible {
			return e.Bid{}, ErrForbidden
		}
	}

	// Create bid
	bid, err := s.bidRepo.Create(ctx, rt.CreateBidInput{
		Name:        in.Name,
		Description: in.Description,
		AuthorType:  in.AuthorType,
		AuthorId:    in.AuthorId,
		TenderId:    in.TenderId,
	})
	if err != nil {
		log.Errorf("BidService - CreateBid - bidRepo.Create: %v", err)
		return e.Bid{}, ErrCreateBid
	}

	return bid, nil
}

func (s *BidService) SubmitDecision(ctx context.Context, bidId uuid.UUID, username string, decision string) (e.Bid, error) {
	// Check if user exists
	user, err := s.employeeRepo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return e.Bid{}, ErrUsername
		}
		log.Errorf("BidService.SubmitDecision - employeeRepo.GetByUsername: %v", err)
		return e.Bid{}, ErrGetEmployeeById
	}

	// Check if bid exists and published
	bid, err := s.bidRepo.Get(ctx, bidId, rt.VersionLatest)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return e.Bid{}, ErrNotFoundBid
		}
		log.Errorf("BidService.SubmitDecision - bidRepo.Get: %v", err)
		return e.Bid{}, ErrGetBid
	}
	if bid.Status != "Published" {
		return e.Bid{}, ErrNotFoundBid
	}

	// Check if tender exists and published
	tender, err := s.tenderRepo.Get(ctx, bid.TenderId, rt.VersionLatest)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return e.Bid{}, ErrNotFoundTender
		}
		log.Errorf("BidService.SubmitDecision - tenderRepo.Get: %v", err)
		return e.Bid{}, ErrGetTender
	}
	if tender.Status != "Published" {
		return e.Bid{}, ErrNotFoundTender
	}

	// Check responsibility
	isResponsible, err := s.employeeRepo.IsResponsible(ctx, tender.OrganizationId, user.Id)
	if err != nil {
		log.Errorf("BidService.SubmitDecision - employeeRepo.IsResponsible: %v", err)
		return e.Bid{}, ErrCheckResponsibility
	}
	if !isResponsible {
		return e.Bid{}, ErrForbidden
	}

	// Make decision
	var resBid e.Bid
	switch decision {
	case "Rejected":
		resBid = bid
	case "Approved":
		if _, err := s.tenderRepo.ChangeStatus(ctx, tender.Id, "Closed"); err != nil {
			if errors.Is(err, repoerrors.ErrNotFound) {
				return e.Bid{}, ErrNotFoundTender
			}
			log.Errorf("BidService.ChangeStatus - tenderRepo.ChangeStatus: %v", err)
			return e.Bid{}, ErrGetTender
		}
		resBid = bid
	}

	return resBid, nil
}

func (s *BidService) ChangeStatus(ctx context.Context, bidId uuid.UUID, status, username string) (e.Bid, error) {
	// Check if user exists
	user, err := s.employeeRepo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return e.Bid{}, ErrUsername
		}
		log.Errorf("BidService.ChangeStatus - employeeRepo.GetByUsername: %v", err)
		return e.Bid{}, ErrGetEmployeeById
	}

	// Check if bid exists
	bid, err := s.bidRepo.Get(ctx, bidId, rt.VersionLatest)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return e.Bid{}, ErrNotFoundBid
		}
		log.Errorf("BidService.ChangeStatus - bidRepo.Get: %v", err)
		return e.Bid{}, ErrGetBid
	}

	// Check responsibility
	var isResponsible bool
	switch bid.AuthorType {
	case "User":
		isResponsible = bid.AuthorId == user.Id
	case "Organization":
		orgId, err := s.employeeRepo.GetOrgIdFromResponsible(ctx, bid.AuthorId)
		if err != nil {
			log.Errorf("BidService.ChangeStatus - employeeRepo.GetOrgIdFromResponsible: %v", err)
			return e.Bid{}, ErrCheckResponsibility
		}
		isResponsible, err = s.employeeRepo.IsResponsible(ctx, orgId, user.Id)
		if err != nil {
			log.Errorf("BidService.ChangeStatus - employeeRepo.IsResponsible: %v", err)
			return e.Bid{}, ErrCheckResponsibility
		}
	}
	if !isResponsible {
		return e.Bid{}, ErrForbidden
	}

	// Change status
	b, err := s.bidRepo.ChangeStatus(ctx, bidId, status)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return e.Bid{}, ErrNotFoundBid
		}
		log.Errorf("BidService.ChangeStatus - bidRepo.ChangeStatus: %v", err)
		return e.Bid{}, ErrGetBid
	}

	return b, nil
}

func (s *BidService) Get(ctx context.Context, bidId uuid.UUID, username string) (e.Bid, error) {
	// Check if user exists
	user, err := s.employeeRepo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return e.Bid{}, ErrUsername
		}
		log.Errorf("BidService.Get - employeeRepo.GetByUsername: %v", err)
		return e.Bid{}, ErrGetEmployeeById
	}

	// Check if bid exists
	bid, err := s.bidRepo.Get(ctx, bidId, rt.VersionLatest)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return e.Bid{}, ErrNotFoundBid
		}
		log.Errorf("BidService.Get - bidRepo.Get: %v", err)
		return e.Bid{}, ErrGetBid
	}

	// Check responsibility
	var isResponsible bool
	switch bid.AuthorType {
	case "User":
		isResponsible = bid.AuthorId == user.Id
	case "Organization":
		orgId, err := s.employeeRepo.GetOrgIdFromResponsible(ctx, bid.AuthorId)
		if err != nil {
			log.Errorf("BidService.Get - employeeRepo.GetOrgIdFromResponsible: %v", err)
			return e.Bid{}, ErrCheckResponsibility
		}
		isResponsible, err = s.employeeRepo.IsResponsible(ctx, orgId, user.Id)
		if err != nil {
			log.Errorf("BidService.Get - employeeRepo.IsResponsible: %v", err)
			return e.Bid{}, ErrCheckResponsibility
		}
	}
	if !isResponsible {
		return e.Bid{}, ErrForbidden
	}

	return bid, nil
}

func (s *BidService) Edit(ctx context.Context, in EditBidInput) (e.Bid, error) {
	// Check if user exists
	user, err := s.employeeRepo.GetByUsername(ctx, in.Username)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return e.Bid{}, ErrUsername
		}
		log.Errorf("BidService.Edit - employeeRepo.GetByUsername: %v", err)
		return e.Bid{}, ErrGetEmployeeById
	}

	// Check if bid exists
	bid, err := s.bidRepo.Get(ctx, in.BidId, rt.VersionLatest)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return e.Bid{}, ErrNotFoundBid
		}
		log.Errorf("BidService.Edit - bidRepo.Get: %v", err)
		return e.Bid{}, ErrGetBid
	}

	// Check responsibility
	var isResponsible bool
	switch bid.AuthorType {
	case "User":
		isResponsible = bid.AuthorId == user.Id
	case "Organization":
		orgId, err := s.employeeRepo.GetOrgIdFromResponsible(ctx, bid.AuthorId)
		if err != nil {
			log.Errorf("BidService.Edit - employeeRepo.GetOrgIdFromResponsible: %v", err)
			return e.Bid{}, ErrCheckResponsibility
		}
		isResponsible, err = s.employeeRepo.IsResponsible(ctx, orgId, user.Id)
		if err != nil {
			log.Errorf("BidService.Edit - employeeRepo.IsResponsible: %v", err)
			return e.Bid{}, ErrCheckResponsibility
		}
	}
	if !isResponsible {
		return e.Bid{}, ErrForbidden
	}

	// Create edited version
	input := rt.CreateSpecifiedBidInput{
		Id:          in.BidId,
		Name:        in.Name,
		Description: in.Description,
		AuthorType:  bid.AuthorType,
		AuthorId:    bid.AuthorId,
		Status:      bid.Status,
		Version:     bid.Version + 1,
		TenderId:    bid.TenderId,
	}
	if in.Name == "" {
		input.Name = bid.Name
	}
	if in.Description == "" {
		input.Description = bid.Description
	}
	b, err := s.bidRepo.CreateSpecified(ctx, input)
	if err != nil {
		log.Errorf("BidService.Edit - bidRepo.CreateSpecified: %v", err)
		return e.Bid{}, ErrCreateBid
	}

	return b, nil
}

func (s *BidService) Rollback(ctx context.Context, bidId uuid.UUID, verison int, username string) (e.Bid, error) {
	// Check if user exists
	user, err := s.employeeRepo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return e.Bid{}, ErrUsername
		}
		log.Errorf("BidService.Rollback - employeeRepo.GetByUsername: %v", err)
		return e.Bid{}, ErrGetEmployeeById
	}

	// Check if bid exists
	bidToRollback, err := s.bidRepo.Get(ctx, bidId, verison)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return e.Bid{}, ErrNotFoundBid
		}
		log.Errorf("BidService.Rollback - bidRepo.Get: %v", err)
		return e.Bid{}, ErrGetBid
	}

	// Check responsibility
	var isResponsible bool
	switch bidToRollback.AuthorType {
	case "User":
		isResponsible = bidToRollback.AuthorId == user.Id
	case "Organization":
		orgId, err := s.employeeRepo.GetOrgIdFromResponsible(ctx, bidToRollback.AuthorId)
		if err != nil {
			log.Errorf("BidService.Rollback - employeeRepo.GetOrgIdFromResponsible: %v", err)
			return e.Bid{}, ErrCheckResponsibility
		}
		isResponsible, err = s.employeeRepo.IsResponsible(ctx, orgId, user.Id)
		if err != nil {
			log.Errorf("BidService.Rollback - employeeRepo.IsResponsible: %v", err)
			return e.Bid{}, ErrCheckResponsibility
		}
	}
	if !isResponsible {
		return e.Bid{}, ErrForbidden
	}

	// Get latest version bid
	latestVersionBid, err := s.bidRepo.Get(ctx, bidId, rt.VersionLatest)
	if err != nil {
		log.Errorf("BidService.Rollback - bidRepo.Get: %v", err)
		return e.Bid{}, ErrGetBid
	}

	// Rollback bid
	b, err := s.bidRepo.CreateSpecified(ctx, rt.CreateSpecifiedBidInput{
		Id: bidId,
		Name: bidToRollback.Name,
		Version: latestVersionBid.Version + 1,
		Description: bidToRollback.Description,
		AuthorType: bidToRollback.AuthorType,
		AuthorId: bidToRollback.AuthorId,
		Status: bidToRollback.Status,
		TenderId: bidToRollback.TenderId,
	})
	if err != nil {
		log.Errorf("BidService.Rollback - bidRepo.CreateSpecified: %v", err)
		return e.Bid{}, ErrCreateBid
	}

	return b, nil
}
