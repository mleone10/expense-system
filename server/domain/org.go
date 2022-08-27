package domain

import (
	"context"
	"errors"
)

var ErrMaxOrgs = errors.New("user has reached the org limit")

type OrgId string

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
}

type OrgRepo interface {
	GetOrg(context.Context, OrgId) (Organization, error)
	GetOrgsForUser(context.Context, UserId) ([]Organization, error)
}
