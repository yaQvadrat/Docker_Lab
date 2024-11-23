package entity

import (
	"time"

	"github.com/google/uuid"
)

type Tender struct {
	Id              uuid.UUID `db:"id"`
	Name            string    `db:"name"`
	Description     string    `db:"description"`
	Type            string    `db:"type"`
	Status          string    `db:"status"`
	OrganizationId  uuid.UUID `db:"organization_id"`
	Version         int       `db:"version"`
	CreatorUsername string    `db:"creator_username"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}
