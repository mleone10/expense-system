package service

import (
	"context"

	"github.com/mleone10/expense-system/domain"
)

type OrgService struct {
	orgRepo domain.OrgRepo
}

func NewOrgService(orgRepo domain.OrgRepo) *OrgService {
	return &OrgService{
		orgRepo: orgRepo,
	}
}

func (s *OrgService) GetOrgsForUser(ctx context.Context, userId domain.UserId) ([]domain.Organization, error) {
	return s.orgRepo.GetOrgsForUser(ctx, userId)
}

func (s *OrgService) CreateOrg(ctx context.Context, name string, adminId domain.UserId) (domain.Organization, error) {
	return domain.Organization{}, nil
}
