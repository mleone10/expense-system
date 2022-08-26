package domain

type OrgId string

type Organization struct {
	Id      OrgId
	Name    string
	Members []User
}

type OrgService interface {
	GetOrgsForUser(UserId) ([]Organization, error)
}

type OrgRepo interface {
	GetOrg(OrgId) (Organization, error)
	GetOrgsForUser(UserId) ([]Organization, error)
}
