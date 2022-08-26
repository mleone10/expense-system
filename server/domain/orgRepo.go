package domain

type OrgRepo interface {
	GetOrg(OrgId) (Organization, error)
	GetOrgsForUser(UserId) ([]Organization, error)
}
