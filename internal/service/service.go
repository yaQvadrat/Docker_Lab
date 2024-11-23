package service

import (
	e "app/internal/entity"
	"app/internal/repo"
	"context"

	"github.com/google/uuid"
)

type CreateTenderInput struct {
	Name            string
	Description     string
	ServiceType     string
	OrganizationId  uuid.UUID
	CreatorUsername string
}

type GetByUsernameInput struct {
	Limit    int
	Offset   int
	Username string
}

type GetTendersInput struct {
	Limit       int
	Offset      int
	ServiceType []string
}

type ChangeTenderStatusInput struct {
	TenderId uuid.UUID
	Status   string
	Username string
}

type EditTenderInput struct {
	TenderId    uuid.UUID
	Username    string
	Name        string
	Description string
	ServiceType string
}

type RollbackTenderInput struct {
	TenderId uuid.UUID
	Version  int
	Username string
}

type Tender interface {
	CreateTender(ctx context.Context, in CreateTenderInput) (e.Tender, error)
	ChangeStatus(ctx context.Context, in ChangeTenderStatusInput) (e.Tender, error)
	Edit(ctx context.Context, in EditTenderInput) (e.Tender, error)
	Rollback(ctx context.Context, in RollbackTenderInput) (e.Tender, error)
	GetTendersByUsername(ctx context.Context, in GetByUsernameInput) ([]e.Tender, error)
	GetTenders(ctx context.Context, in GetTendersInput) ([]e.Tender, error)
	GetTender(ctx context.Context, tenderId uuid.UUID, username string) (e.Tender, error)
}

type CreateBidInput struct {
	Name        string
	Description string
	TenderId    uuid.UUID
	AuthorType  string
	AuthorId    uuid.UUID
}

type EditBidInput struct {
	BidId       uuid.UUID
	Username    string
	Name        string
	Description string
}

type Bid interface {
	CreateBid(ctx context.Context, in CreateBidInput) (e.Bid, error)
	SubmitDecision(ctx context.Context, bidId uuid.UUID, username string, decision string) (e.Bid, error)
	ChangeStatus(ctx context.Context, bidId uuid.UUID, status, username string) (e.Bid, error)
	Get(ctx context.Context, bidId uuid.UUID, username string) (e.Bid, error)
	Edit(ctx context.Context, in EditBidInput) (e.Bid, error)
	Rollback(ctx context.Context, bidId uuid.UUID, version int, username string) (e.Bid, error)
}

type Services struct {
	Tender
	Bid
}

type ServicesDependencies struct {
	Repos *repo.Repositories
}

func NewServices(d ServicesDependencies) *Services {
	return &Services{
		Tender: NewTenderService(d.Repos.Tender, d.Repos.Employee),
		Bid:    NewBidService(d.Repos.Tender, d.Repos.Employee, d.Repos.Bid),
	}
}
