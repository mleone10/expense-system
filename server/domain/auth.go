package domain

type AuthClient interface {
	GetAuthTokens(string) (AuthTokens, error)
	RedirectUrl() string
}

// AuthTokens is a container for the tokens returned from the Auth Code flow.
type AuthTokens struct {
	AccessToken string
}
