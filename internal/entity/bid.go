package entity

import (
	"time"

	"github.com/google/uuid"
)

type Bid struct {
	Id          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	AuthorType  string    `db:"author"`
	AuthorId    uuid.UUID `db:"author_id"`
	Status      string    `db:"status"`
	Version     int       `db:"version"`
	TenderId    uuid.UUID `db:"tender_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
