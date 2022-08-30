package service

import (
	"context"

	"github.com/mleone10/expense-system/domain"
)

type AuthenticatedUserService struct{}

func NewAuthedUserService() *AuthenticatedUserService {
	return &AuthenticatedUserService{}
}

func (s *AuthenticatedUserService) GetAuthenticatedUserInfo(ctx context.Context, authToken string) (domain.User, error) {
	return domain.User{}, nil
}
