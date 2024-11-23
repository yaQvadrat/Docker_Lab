package pgdb

import (
	e "app/internal/entity"
	"app/internal/repo/repoerrors"
	rt "app/internal/repo/repotypes"
	"app/pkg/postgres"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type BidRepo struct {
	*postgres.Postgres
}

func NewBidRepo(pg *postgres.Postgres) *BidRepo {
	return &BidRepo{pg}
}

// use repotypes.VersionLatest if need latest version
func (r *BidRepo) Get(ctx context.Context, id uuid.UUID, version int) (e.Bid, error) {
	sql := `
		SELECT DISTINCT ON (id) * FROM bid
		WHERE id = $1
	`
	if version != rt.VersionLatest {
		sql = fmt.Sprintf("%s AND version = %d", sql, version)
	}
	sql += "\n\t\tORDER BY id, version DESC"

	rows, err := r.Pool.Query(ctx, sql, id)
	if err != nil {
		return e.Bid{}, fmt.Errorf("pgdb - BidRepo.Get - Pool.Query: %w", err)
	}

	b, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[e.Bid])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return e.Bid{}, repoerrors.ErrNotFound
		}
		return e.Bid{}, fmt.Errorf("pgdb - BidRepo.Get - CollectExactlyOneRow: %w", err)
	}

	return b, nil
}

func (r *BidRepo) Create(ctx context.Context, in rt.CreateBidInput) (e.Bid, error) {
	sql := `
		INSERT INTO bid
			(name, description, author, author_id, tender_id)
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING *
	`

	rows, err := r.Pool.Query(ctx, sql,
		in.Name,
		in.Description,
		in.AuthorType,
		in.AuthorId,
		in.TenderId,
	)
	if err != nil {
		return e.Bid{}, fmt.Errorf("pgdb.BidRepo - Create - Pool.Query: %w", err)
	}

	b, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[e.Bid])
	if err != nil {
		return e.Bid{}, fmt.Errorf("pgdb.BidRepo - Create - pgx.CollectExactlyOneRow: %w", err)
	}

	return b, nil
}

func (r *BidRepo) ChangeStatus(ctx context.Context, id uuid.UUID, status string) (e.Bid, error) {
	sql := `
		UPDATE bid
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE
			id = $2
  			AND version = (
      			SELECT MAX(version)
      			FROM bid
      			WHERE id = $2
  			)
		RETURNING *;
	`

	rows, err := r.Pool.Query(ctx, sql, status, id)
	if err != nil {
		return e.Bid{}, fmt.Errorf("pgdb - BidRepo.ChangeStatus - Pool.Query: %w", err)
	}

	b, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[e.Bid])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return e.Bid{}, repoerrors.ErrNotFound
		}
		return e.Bid{}, fmt.Errorf("pgdb - BidRepo.ChangeStatus - CollectExactlyOneRow: %w", err)
	}

	return b, nil
}

func (r *BidRepo) CreateSpecified(ctx context.Context, in rt.CreateSpecifiedBidInput) (e.Bid, error) {
	sql := `
		INSERT INTO bid
			(id, name, description, author, author_id, status, version, tender_id)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING *
	`

	rows, err := r.Pool.Query(ctx, sql,
		in.Id,
		in.Name,
		in.Description,
		in.AuthorType,
		in.AuthorId,
		in.Status,
		in.Version,
		in.TenderId,
	)
	if err != nil {
		return e.Bid{}, fmt.Errorf("pgdb - BidRepo.CreateSpecified - Pool.Query: %w", err)
	}

	b, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[e.Bid])
	if err != nil {
		return e.Bid{}, fmt.Errorf("pgdb - BidRepo.CreateSpecified - pgx.CollectExactlyOneRow: %w", err)
	}

	return b, nil
}
