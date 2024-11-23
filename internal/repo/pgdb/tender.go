package pgdb

import (
	e "app/internal/entity"
	"app/internal/repo/repoerrors"
	rt "app/internal/repo/repotypes"
	"app/pkg/postgres"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type TenderRepo struct {
	*postgres.Postgres
}

func NewTenderRepo(pg *postgres.Postgres) *TenderRepo {
	return &TenderRepo{pg}
}

func (r *TenderRepo) CreateTender(ctx context.Context, in rt.CreateTenderInput) (e.Tender, error) {
	sql := `
		INSERT INTO tender
			(name, description, type, organization_id, creator_username)
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING *
	`

	rows, err := r.Pool.Query(ctx, sql,
		in.Name,
		in.Description,
		in.ServiceType,
		in.OrganizationId,
		in.CreatorUsername,
	)
	if err != nil {
		return e.Tender{}, fmt.Errorf("pgdb - CreateTender - Pool.Query: %w", err)
	}

	t, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[e.Tender])
	if err != nil {
		return e.Tender{}, fmt.Errorf("pgdb - CreateTender - pgx.CollectExactlyOneRow: %w", err)
	}

	return t, nil
}

func (r *TenderRepo) GetTendersByUsername(ctx context.Context, in rt.GetByUsernameInput) ([]e.Tender, error) {
	sql := `
		SELECT * 
		FROM (
			SELECT DISTINCT ON (id) * FROM tender
			WHERE creator_username = $1
			ORDER BY id, version DESC
		) AS last_versions
		ORDER BY name
		LIMIT $2 OFFSET $3
	`

	rows, err := r.Pool.Query(ctx, sql, in.Username, in.Limit, in.Offset)
	if err != nil {
		return nil, fmt.Errorf("pgdb - GetTendersByUsername - Pool.Query: %w", err)
	}

	tenders, err := pgx.CollectRows(rows, pgx.RowToStructByName[e.Tender])
	if err != nil {
		return nil, fmt.Errorf("pgdb - GetTendersByUsername - CollectRows: %w", err)
	}

	return tenders, nil
}

func (r *TenderRepo) GetPublishedTenders(ctx context.Context, in rt.GetPublishedTendersInput) ([]e.Tender, error) {
	sql := `
		SELECT *
		FROM (
			SELECT DISTINCT ON (id) * FROM tender
			ORDER BY id, version DESC
		) AS last_versions
		WHERE status='Published'
	`
	if in.ServiceType != nil {
		for i, s := range in.ServiceType {
			in.ServiceType[i] = fmt.Sprintf("'%s'", s)
		}
		sql = fmt.Sprintf("%s\tAND type IN (%s)", sql, strings.Join(in.ServiceType, ","))
	}
	sql += "\n\t\tORDER BY name\n\t\tLIMIT $1 OFFSET $2"

	rows, err := r.Pool.Query(ctx, sql, in.Limit, in.Offset)
	if err != nil {
		return nil, fmt.Errorf("pgdb - GetPublishedTenders - Pool.Query: %w", err)
	}

	tenders, err := pgx.CollectRows(rows, pgx.RowToStructByName[e.Tender])
	if err != nil {
		return nil, fmt.Errorf("pgdb - GetPublishedTenders - CollectRows: %w", err)
	}

	return tenders, nil
}

// use repotypes.VersionLatest if need latest version
func (r *TenderRepo) Get(ctx context.Context, id uuid.UUID, version int) (e.Tender, error) {
	sql := `
		SELECT DISTINCT ON (id) * FROM tender
		WHERE id = $1
	`
	if version != rt.VersionLatest {
		sql = fmt.Sprintf("%s AND version = %d", sql, version)
	}
	sql += "\n\t\tORDER BY id, version DESC"

	rows, err := r.Pool.Query(ctx, sql, id)
	if err != nil {
		return e.Tender{}, fmt.Errorf("pgdb - TenderRepo.Get - Pool.Query: %w", err)
	}

	t, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[e.Tender])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return e.Tender{}, repoerrors.ErrNotFound
		}
		return e.Tender{}, fmt.Errorf("pgdb - TenderRepo.Get - CollectExactlyOneRow: %w", err)
	}

	return t, nil
}

func (r *TenderRepo) ChangeStatus(ctx context.Context, id uuid.UUID, status string) (e.Tender, error) {
	sql := `
		UPDATE tender
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE
			id = $2
  			AND version = (
      			SELECT MAX(version)
      			FROM tender
      			WHERE id = $2
  			)
		RETURNING *;
	`

	rows, err := r.Pool.Query(ctx, sql, status, id)
	if err != nil {
		return e.Tender{}, fmt.Errorf("pgdb - TenderRepo.ChangeStatus - Pool.Query: %w", err)
	}

	t, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[e.Tender])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return e.Tender{}, repoerrors.ErrNotFound
		}
		return e.Tender{}, fmt.Errorf("pgdb - TenderRepo.ChangeStatus - CollectExactlyOneRow: %w", err)
	}

	return t, nil
}

func (r *TenderRepo) CreateSpecified(ctx context.Context, in rt.CreateSpecifiedInput) (e.Tender, error) {
	sql := `
		INSERT INTO tender
			(id, name, description, type, organization_id, version, creator_username, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING *
	`

	rows, err := r.Pool.Query(ctx, sql,
		in.Id,
		in.Name,
		in.Description,
		in.ServiceType,
		in.OrganizationId,
		in.Version,
		in.CreatorUsername,
		in.Status,
	)
	if err != nil {
		return e.Tender{}, fmt.Errorf("pgdb - TenderRepo.CreateSpecified - Pool.Query: %w", err)
	}

	t, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[e.Tender])
	if err != nil {
		return e.Tender{}, fmt.Errorf("pgdb - TenderRepo.CreateSpecified - pgx.CollectExactlyOneRow: %w", err)
	}

	return t, nil
}

func (r *TenderRepo) GetLatestVersion(ctx context.Context, id uuid.UUID) (int, error) {
	sql := `
		SELECT DISTINCT ON (id) version FROM tender
		WHERE id = $1
		ORDER BY id, version DESC
	`

	rows, err := r.Pool.Query(ctx, sql, id)
	if err != nil {
		return 0, fmt.Errorf("pgdb - TenderRepo.GetLatestVersion - Pool.Query: %w", err)
	}

	latestVersion, err := pgx.CollectExactlyOneRow(rows, pgx.RowTo[int])
	if err != nil {
		return 0, fmt.Errorf("pgdb - TenderRepo.GetLatestVersion - pgx.CollectExactlyOneRow: %w", err)
	}

	return latestVersion, nil
}
