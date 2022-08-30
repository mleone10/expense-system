package service

import (
	"context"

	"github.com/mleone10/expense-system/domain"
)

type AuthenticatedUserService struct {
	authenticatedUserRepo domain.AuthenticatedUserRepo
}

func NewAuthenticatedUserService(repo domain.AuthenticatedUserRepo) *AuthenticatedUserService {
	return &AuthenticatedUserService{
		authenticatedUserRepo: repo,
	}
}

func (s *AuthenticatedUserService) GetAuthenticatedUserInfo(ctx context.Context, authToken string) (domain.User, error) {
	return s.authenticatedUserRepo.GetAuthenticatedUserInfo(ctx, authToken)
}
