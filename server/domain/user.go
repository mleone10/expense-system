package domain

import "context"

type UserId string

type User struct {
	Id         UserId
	Name       string
	ProfileUrl string
}

type AuthenticatedUserService interface {
	GetAuthenticatedUserInfo(ctx context.Context, authToken string) (User, error)
}

type AuthenticatedUserRepo interface {
	GetAuthenticatedUserInfo(ctx context.Context, authToken string) (User, error)
}
