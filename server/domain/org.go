package domain

import "context"

type OrgId string

type Organization struct {
	Id      OrgId
	Name    string
	Members []User
}

type OrgService interface {
	GetOrgsForUser(context.Context, UserId) ([]Organization, error)
	CreateOrg(ctx context.Context, name string, adminId UserId) (Organization, error)
}

type OrgRepo interface {
	GetOrg(context.Context, OrgId) (Organization, error)
	GetOrgsForUser(context.Context, UserId) ([]Organization, error)
}
