package cognito

import "net/http"

// AuthClient is a container for various auth-related configuration and serves as a receiver for all auth functionality.
type CognitoClient struct {
	client              http.Client
	clientHostname      string
	clientScheme        string
	cognitoClientId     string
	cognitoClientSecret string
}

// OptionFuncs are used to configure an AuthClient using Functional Configuration
type OptionFunc func(*CognitoClient)

// NewAuthClient applies the given OptionFuncs to a new AuthClient before returning it, ready to use.
func NewAuthClient(options ...OptionFunc) *CognitoClient {
	a := &CognitoClient{}

	for _, opt := range options {
		opt(a)
	}

	a.client = http.Client{}

	return a
}

// WithClientHostname is an OptionFunc used to set the hostname of the current client.
func WithClientHostname(hostname string) OptionFunc {
	return func(ac *CognitoClient) {
		ac.clientHostname = hostname
	}
}

// WithClientScheme should usually be "https" unless running locally, in which case "http" is used instead.
func WithClientScheme(scheme string) OptionFunc {
	return func(ac *CognitoClient) {
		ac.clientScheme = scheme
	}
}

// WithCognitoClientId sets the public client ID used to authenticate with AWS Cognito.
func WithCognitoClientId(clientId string) OptionFunc {
	return func(ac *CognitoClient) {
		ac.cognitoClientId = clientId
	}
}

// WithCognitoClientSecret sets the private client secret used to authenticate with AWS Cognito.
func WithCognitoClientSecret(clientSecret string) OptionFunc {
	return func(ac *CognitoClient) {
		ac.cognitoClientSecret = clientSecret
	}
}
