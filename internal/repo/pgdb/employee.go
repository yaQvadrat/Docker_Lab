package pgdb

import (
	e "app/internal/entity"
	"app/internal/repo/repoerrors"
	"app/pkg/postgres"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type EmployeeRepo struct {
	*postgres.Postgres
}

func NewEmployeeRepo(pg *postgres.Postgres) *EmployeeRepo {
	return &EmployeeRepo{pg}
}

func (r *EmployeeRepo) GetByUsername(ctx context.Context, username string) (e.Employee, error) {
	sql := `
		SELECT * FROM employee
		WHERE username = $1
	`

	rows, err := r.Pool.Query(ctx, sql, username)
	if err != nil {
		return e.Employee{}, fmt.Errorf("pgdb - GetByUsername - Query: %w", err)
	}

	employee, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[e.Employee])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return e.Employee{}, repoerrors.ErrNotFound
		}
		return e.Employee{}, fmt.Errorf("pgdb - GetByUsername - CollectExactlyOneRow: %w", err)
	}

	return employee, nil
}

func (r *EmployeeRepo) IsResponsible(ctx context.Context, orgId, userId uuid.UUID) (bool, error) {
	sql := `
		SELECT EXISTS(
    		SELECT 1
    		FROM organization_responsible
    		WHERE organization_id = $1
    		AND user_id = $2
		) AS is_responsible
	`

	var isResponsible bool
	err := r.Pool.QueryRow(ctx, sql, orgId, userId).Scan(&isResponsible)
	if err != nil {
		return false, fmt.Errorf("pgdb - IsResponsible - QueryRow: %w", err)
	}

	return isResponsible, nil
}

func (r *EmployeeRepo) IsResponsibleSimplified(ctx context.Context, userId uuid.UUID) (bool, error) {
	sql := `
		SELECT EXISTS(
    		SELECT 1
    		FROM organization_responsible
    		WHERE user_id = $1
		) AS is_responsible
	`

	var isResponsible bool
	err := r.Pool.QueryRow(ctx, sql, userId).Scan(&isResponsible)
	if err != nil {
		return false, fmt.Errorf("pgdb - IsResponsibleSimplified - QueryRow: %w", err)
	}

	return isResponsible, nil
}

func (r *EmployeeRepo) GetOrgIdFromResponsible(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	sql := `
		SELECT organization_id
		FROM organization_responsible
		WHERE user_id = $1
	`

	var orgId uuid.UUID
	err := r.Pool.QueryRow(ctx, sql, id).Scan(&orgId)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("pgdb - GetOrgIdFromResponsible - QueryRow: %w", err)
	}

	return orgId, nil
}

func (r *EmployeeRepo) GetById(ctx context.Context, id uuid.UUID) (e.Employee, error) {
	sql := `
		SELECT * FROM employee
		WHERE id = $1
	`

	rows, err := r.Pool.Query(ctx, sql, id)
	if err != nil {
		return e.Employee{}, fmt.Errorf("pgdb - GetById - Query: %w", err)
	}

	employee, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[e.Employee])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return e.Employee{}, repoerrors.ErrNotFound
		}
		return e.Employee{}, fmt.Errorf("pgdb - GetById - CollectExactlyOneRow: %w", err)
	}

	return employee, nil
}
