package repotypes

import "github.com/google/uuid"

const (
	VersionLatest = 0
)

type CreateTenderInput struct {
	Name            string
	Description     string
	ServiceType     string
	OrganizationId  uuid.UUID
	CreatorUsername string
	Status          string
}

type GetByUsernameInput struct {
	Limit    int
	Offset   int
	Username string
}

type GetPublishedTendersInput struct {
	Limit       int
	Offset      int
	ServiceType []string
}

type CreateSpecifiedInput struct {
	Id      uuid.UUID
	Version int
	CreateTenderInput
}
