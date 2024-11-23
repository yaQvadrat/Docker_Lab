package repotypes

import "github.com/google/uuid"

type CreateBidInput struct {
	Name        string
	Description string
	TenderId    uuid.UUID
	AuthorType  string
	AuthorId    uuid.UUID
}

type CreateSpecifiedBidInput struct {
	Id          uuid.UUID
	Name        string
	Description string
	AuthorType  string
	AuthorId    uuid.UUID
	Status      string
	Version     int
	TenderId    uuid.UUID
}
