package service

import (
	"context"
	"fmt"

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
	if authToken == "" {
		return domain.User{}, fmt.Errorf("cannot retrieve user info with empty auth token")
	}
	return s.authenticatedUserRepo.GetAuthenticatedUserInfo(ctx, authToken)
}
