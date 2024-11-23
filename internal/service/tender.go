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

type TenderService struct {
	tenderRepo   repo.Tender
	employeeRepo repo.Employee
}

func NewTenderService(tRepo repo.Tender, eRepo repo.Employee) *TenderService {
	return &TenderService{
		tenderRepo:   tRepo,
		employeeRepo: eRepo,
	}
}

func (s *TenderService) CreateTender(ctx context.Context, in CreateTenderInput) (e.Tender, error) {
	user, err := s.employeeRepo.GetByUsername(ctx, in.CreatorUsername)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return e.Tender{}, ErrUsername
		}
		log.Errorf("TenderService.CreateTender - employeeRepo.GetByUsername: %v", err)
		return e.Tender{}, ErrGetEmployeeByUsername
	}

	isResponsible, err := s.employeeRepo.IsResponsible(ctx, in.OrganizationId, user.Id)
	if err != nil {
		log.Errorf("TenderService.CreateTender - employeeRepo.IsResponsible: %v", err)
		return e.Tender{}, ErrCheckResponsibility
	}
	if !isResponsible {
		return e.Tender{}, ErrForbidden
	}

	tender, err := s.tenderRepo.CreateTender(ctx, rt.CreateTenderInput{
		Name:            in.Name,
		Description:     in.Description,
		ServiceType:     in.ServiceType,
		OrganizationId:  in.OrganizationId,
		CreatorUsername: in.CreatorUsername,
	})
	if err != nil {
		log.Errorf("TenderService.CreateTender - tenderRepo.CreateTender: %v", err)
		return e.Tender{}, ErrCreateTender
	}

	return tender, nil
}

func (s *TenderService) GetTendersByUsername(ctx context.Context, in GetByUsernameInput) ([]e.Tender, error) {
	_, err := s.employeeRepo.GetByUsername(ctx, in.Username)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return nil, ErrUsername
		}
		log.Errorf("TenderService.GetTendersByUsername - employeeRepo.GetByUsername: %v", err)
		return nil, ErrGetEmployeeByUsername
	}

	tenders, err := s.tenderRepo.GetTendersByUsername(ctx, rt.GetByUsernameInput{
		Limit:    in.Limit,
		Offset:   in.Offset,
		Username: in.Username,
	})
	if err != nil {
		log.Errorf("TenderService.GetTendersByUsername - tenderRepo.GetTendersByUsername: %v", err)
		return nil, ErrGetTendersByUsername
	}

	return tenders, nil
}

func (s *TenderService) GetTenders(ctx context.Context, in GetTendersInput) ([]e.Tender, error) {
	tenders, err := s.tenderRepo.GetPublishedTenders(ctx, rt.GetPublishedTendersInput{
		Limit:       in.Limit,
		Offset:      in.Offset,
		ServiceType: in.ServiceType,
	})
	if err != nil {
		log.Errorf("TenderService.GetTenders - tenderRepo.GetPublishedTenders: %v", err)
		return nil, ErrGetTenders
	}

	return tenders, nil
}

func (s *TenderService) ChangeStatus(ctx context.Context, in ChangeTenderStatusInput) (e.Tender, error) {
	// Check if user exists
	user, err := s.employeeRepo.GetByUsername(ctx, in.Username)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return e.Tender{}, ErrUsername
		}
		log.Errorf("TenderService.ChangeStatus - employeeRepo.GetByUsername: %v", err)
		return e.Tender{}, ErrGetEmployeeByUsername
	}

	// Check if tender exists
	tender, err := s.tenderRepo.Get(ctx, in.TenderId, rt.VersionLatest)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return e.Tender{}, ErrNotFoundTender
		}
		log.Errorf("TenderService.ChangeStatus - tenderRepo.Get: %v", err)
		return e.Tender{}, ErrGetTender
	}

	// Check rights
	isResponsible, err := s.employeeRepo.IsResponsible(ctx, tender.OrganizationId, user.Id)
	if err != nil {
		log.Errorf("TenderService.ChangeStatus - employeeRepo.IsResponsible: %v", err)
		return e.Tender{}, ErrCheckResponsibility
	}
	if !isResponsible {
		return e.Tender{}, ErrForbidden
	}

	// Change status
	t, err := s.tenderRepo.ChangeStatus(ctx, in.TenderId, in.Status)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return e.Tender{}, ErrNotFoundTender
		}
		log.Errorf("TenderService.ChangeStatus - tenderRepo.ChangeStatus: %v", err)
		return e.Tender{}, ErrGetTender
	}

	return t, nil
}

func (s *TenderService) Edit(ctx context.Context, in EditTenderInput) (e.Tender, error) {
	// Check if user exists
	user, err := s.employeeRepo.GetByUsername(ctx, in.Username)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return e.Tender{}, ErrUsername
		}
		log.Errorf("TenderService.Edit - employeeRepo.GetByUsername: %v", err)
		return e.Tender{}, ErrGetEmployeeByUsername
	}

	// Check if tender exists
	tender, err := s.tenderRepo.Get(ctx, in.TenderId, rt.VersionLatest)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return e.Tender{}, ErrNotFoundTender
		}
		log.Errorf("TenderService.Edit - tenderRepo.Get: %v", err)
		return e.Tender{}, ErrGetTender
	}

	// Check rights
	isResponsible, err := s.employeeRepo.IsResponsible(ctx, tender.OrganizationId, user.Id)
	if err != nil {
		log.Errorf("TenderService.Edit - employeeRepo.IsResponsible: %v", err)
		return e.Tender{}, ErrCheckResponsibility
	}
	if !isResponsible {
		return e.Tender{}, ErrForbidden
	}

	// Create edited version
	input := rt.CreateSpecifiedInput{
		Id:      in.TenderId,
		Version: tender.Version + 1,
		CreateTenderInput: rt.CreateTenderInput{
			Name:            in.Name,
			Description:     in.Description,
			ServiceType:     in.ServiceType,
			OrganizationId:  tender.OrganizationId,
			CreatorUsername: tender.CreatorUsername,
			Status:          tender.Status,
		},
	}
	if in.Name == "" {
		input.Name = tender.Name
	}
	if in.Description == "" {
		input.Description = tender.Description
	}
	if in.ServiceType == "" {
		input.ServiceType = tender.Type
	}
	t, err := s.tenderRepo.CreateSpecified(ctx, input)
	if err != nil {
		log.Errorf("TenderService.Edit - tenderRepo.CreateSpecified: %v", err)
		return e.Tender{}, ErrCreateTender
	}

	return t, nil
}

func (s *TenderService) Rollback(ctx context.Context, in RollbackTenderInput) (e.Tender, error) {
	// Check if user exists
	user, err := s.employeeRepo.GetByUsername(ctx, in.Username)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return e.Tender{}, ErrUsername
		}
		log.Errorf("TenderService.Rollback - employeeRepo.GetByUsername: %v", err)
		return e.Tender{}, ErrGetEmployeeByUsername
	}

	// Check if tender exists
	tenderToRollback, err := s.tenderRepo.Get(ctx, in.TenderId, in.Version)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return e.Tender{}, ErrNotFoundTender
		}
		log.Errorf("TenderService.Rollback - tenderRepo.Get: %v", err)
		return e.Tender{}, ErrGetTender
	}

	// Check rights
	isResponsible, err := s.employeeRepo.IsResponsible(ctx, tenderToRollback.OrganizationId, user.Id)
	if err != nil {
		log.Errorf("TenderService.Rollback - employeeRepo.IsResponsible: %v", err)
		return e.Tender{}, ErrCheckResponsibility
	}
	if !isResponsible {
		return e.Tender{}, ErrForbidden
	}

	// Create rollback version
	latestVersion, err := s.tenderRepo.GetLatestVersion(ctx, in.TenderId)
	if err != nil {
		log.Errorf("TenderService.Rollback - tenderRepo.GetLatestVersion: %v", err)
		return e.Tender{}, ErrGetTenderLatestVersion
	}
	t, err := s.tenderRepo.CreateSpecified(ctx, rt.CreateSpecifiedInput{
		Id:      in.TenderId,
		Version: latestVersion + 1,
		CreateTenderInput: rt.CreateTenderInput{
			Name:            tenderToRollback.Name,
			Description:     tenderToRollback.Description,
			ServiceType:     tenderToRollback.Type,
			OrganizationId:  tenderToRollback.OrganizationId,
			CreatorUsername: tenderToRollback.CreatorUsername,
			Status:          tenderToRollback.Status,
		},
	})
	if err != nil {
		log.Errorf("TenderService.Rollback - tenderRepo.CreateSpecified: %v", err)
		return e.Tender{}, ErrCreateTender
	}

	return t, nil
}

func (s *TenderService) GetTender(ctx context.Context, tenderId uuid.UUID, username string) (e.Tender, error) {
	// Check if user exists
	user, err := s.employeeRepo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return e.Tender{}, ErrUsername
		}
		log.Errorf("TenderService.GetTender - employeeRepo.GetByUsername: %v", err)
		return e.Tender{}, ErrGetEmployeeByUsername
	}

	// Check if tender exists
	tender, err := s.tenderRepo.Get(ctx, tenderId, rt.VersionLatest)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return e.Tender{}, ErrNotFoundTender
		}
		log.Errorf("TenderService.GetTender - tenderRepo.Get: %v", err)
		return e.Tender{}, ErrGetTender
	}

	// Check rights
	isResponsible, err := s.employeeRepo.IsResponsible(ctx, tender.OrganizationId, user.Id)
	if err != nil {
		log.Errorf("TenderService.GetTender - employeeRepo.IsResponsible: %v", err)
		return e.Tender{}, ErrCheckResponsibility
	}
	if !isResponsible {
		return e.Tender{}, ErrForbidden
	}

	return tender, nil
}
