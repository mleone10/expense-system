package domain

import (
	"context"
	"errors"

	"github.com/gofrs/uuid"
)

var ErrInvalidRequest = errors.New("request failed validation")

type OrgId string

func NewOrgId() (OrgId, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return OrgId(id.String()), nil
}

type Organization struct {
	Id      OrgId
	Name    string
	Members []Member
}

type Member struct {
	Id    UserId
	Admin bool
}

func (o Organization) IsAdmin(id UserId) bool {
	for _, m := range o.Members {
		if m.Id == id && m.Admin {
			return true
		}
	}
	return false
}

type OrgService interface {
	GetOrgsForUser(context.Context, UserId) ([]Organization, error)
	CreateOrg(ctx context.Context, name string, adminId UserId) (Organization, error)
	GetOrg(ctx context.Context, userId UserId, orgId OrgId) (Organization, error)
}

type OrgRepo interface {
	GetOrg(context.Context, OrgId) (Organization, error)
	GetOrgsForUser(context.Context, UserId) ([]Organization, error)
	CreateOrg(ctx context.Context, name string, adminId UserId) (OrgId, error)
}
