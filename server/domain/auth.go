package domain

import "github.com/lestrrat-go/jwx/jwt"

// An AuthClient is capable of performing authentication via the OAuth Auth Code flow.
type AuthClient interface {
	GetAuthTokens(string) (AuthTokens, error)
	TokenIsValid(string) (jwt.Token, error)
	RedirectUrl() string
}

// AuthTokens is a container for the tokens returned from the Auth Code flow.
type AuthTokens struct {
	AccessToken string
}
