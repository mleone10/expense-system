package service

import (
	"context"
	"fmt"

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
	if name == "" {
		return domain.Organization{}, domain.ErrInvalidRequest
	}

	orgs, err := s.orgRepo.GetOrgsForUser(ctx, adminId)
	if err != nil {
		return domain.Organization{}, fmt.Errorf("failed to list existing orgs: %w", err)
	}

	if numOrgsAsAdmin(orgs, adminId) >= 3 {
		return domain.Organization{}, domain.ErrInvalidRequest
	}

	orgId, err := s.orgRepo.CreateOrg(ctx, name, adminId)
	if err != nil {
		return domain.Organization{}, fmt.Errorf("failed to create org: %w", err)
	}

	return domain.Organization{
		Id:   orgId,
		Name: name,
		Members: []domain.Member{
			{
				Id:    adminId,
				Admin: true,
			},
		},
	}, nil
}

func (s *OrgService) GetOrg(ctx context.Context, userId domain.UserId, orgId domain.OrgId) (domain.Organization, error) {
	return domain.Organization{
		Id:      orgId,
		Name:    "",
		Members: []domain.Member{{Id: userId, Admin: false}},
	}, nil
}

func numOrgsAsAdmin(orgs []domain.Organization, userId domain.UserId) int {
	adminCount := 0
	for _, o := range orgs {
		if o.IsAdmin(userId) {
			adminCount++
		}
	}
	return adminCount
}
