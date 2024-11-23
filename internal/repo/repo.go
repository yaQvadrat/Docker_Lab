package repo

import (
	e "app/internal/entity"
	"app/internal/repo/pgdb"
	rt "app/internal/repo/repotypes"
	"app/pkg/postgres"
	"context"

	"github.com/google/uuid"
)

type Tender interface {
	CreateTender(ctx context.Context, in rt.CreateTenderInput) (e.Tender, error)
	ChangeStatus(ctx context.Context, id uuid.UUID, status string) (e.Tender, error)
	CreateSpecified(ctx context.Context, in rt.CreateSpecifiedInput) (e.Tender, error)
	Get(ctx context.Context, id uuid.UUID, version int) (e.Tender, error)
	GetTendersByUsername(ctx context.Context, in rt.GetByUsernameInput) ([]e.Tender, error)
	GetPublishedTenders(ctx context.Context, in rt.GetPublishedTendersInput) ([]e.Tender, error)
	GetLatestVersion(ctx context.Context, id uuid.UUID) (int, error)
}

type Employee interface {
	IsResponsible(ctx context.Context, orgId, userId uuid.UUID) (bool, error)
	IsResponsibleSimplified(ctx context.Context, userId uuid.UUID) (bool, error)
	GetByUsername(ctx context.Context, username string) (e.Employee, error)
	GetById(ctx context.Context, id uuid.UUID) (e.Employee, error)
	GetOrgIdFromResponsible(ctx context.Context, id uuid.UUID) (uuid.UUID, error)
}

type Bid interface {
	Get(ctx context.Context, id uuid.UUID, version int) (e.Bid, error)
	Create(ctx context.Context, in rt.CreateBidInput) (e.Bid, error)
	CreateSpecified(ctx context.Context, in rt.CreateSpecifiedBidInput) (e.Bid, error)
	ChangeStatus(ctx context.Context, id uuid.UUID, status string) (e.Bid, error)
}

type Repositories struct {
	Tender
	Employee
	Bid
}

func NewPostgresRepo(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		Tender:   pgdb.NewTenderRepo(pg),
		Employee: pgdb.NewEmployeeRepo(pg),
		Bid:      pgdb.NewBidRepo(pg),
	}
}
