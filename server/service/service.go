package service

import "github.com/mleone10/expense-system/domain"

type OrgService struct {
	orgRepo domain.OrgRepo
}

func NewOrgService(orgRepo domain.OrgRepo) *OrgService {
	return &OrgService{
		orgRepo: orgRepo,
	}
}

func (s *OrgService) GetOrgsForUser(userId domain.UserId) ([]domain.Organization, error) {
	return s.orgRepo.GetOrgsForUser(userId)
}
