package googleauth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/mleone10/expense-system/domain"
)

// AuthClient is a container for various auth-related configuration and serves as a receiver for all auth functionality.
type AuthClient struct {
	client              http.Client
	clientHostname      string
	clientScheme        string
	cognitoClientId     string
	cognitoClientSecret string
	basicAuth           string
	authCodeRedirectUri string
}

// OptionFuncs are used to configure an AuthClient using Functional Configuration
type OptionFunc func(*AuthClient)

// NewAuthClient applies the given OptionFuncs to a new AuthClient before returning it, ready to use.
func NewAuthClient(options ...OptionFunc) *AuthClient {
	a := &AuthClient{}

	for _, opt := range options {
		opt(a)
	}

	a.client = http.Client{}

	return a
}

// WithClientHostname is an OptionFunc used to set the hostname of the current client.
func WithClientHostname(hostname string) OptionFunc {
	return func(ac *AuthClient) {
		ac.clientHostname = hostname
	}
}

// WithClientScheme should usually be "https" unless running locally, in which case "http" is used instead.
func WithClientScheme(scheme string) OptionFunc {
	return func(ac *AuthClient) {
		ac.clientScheme = scheme
	}
}

// WithCognitoClientId sets the public client ID used to authenticate with AWS Cognito.
func WithCognitoClientId(clientId string) OptionFunc {
	return func(ac *AuthClient) {
		ac.cognitoClientId = clientId
	}
}

// WithCognitoClientSecret sets the private client secret used to authenticate with AWS Cognito.
func WithCognitoClientSecret(clientSecret string) OptionFunc {
	return func(ac *AuthClient) {
		ac.cognitoClientSecret = clientSecret
	}
}

// RedirectUrl returns the formatted client URL indicated by the previously set client scheme and hostname.  This is where user are redirected after authentication.
func (a *AuthClient) RedirectUrl() string {
	return fmt.Sprintf("%s://%s", a.clientScheme, a.clientHostname)
}

// GetAuthTokens is used during the Auth Code flow to exchange an authentication code for valid OAuth Tokens.
func (a *AuthClient) GetAuthTokens(authCode string) (domain.AuthTokens, error) {
	type tokenResponse struct {
		Id      string `json:"id_token"`
		Access  string `json:"access_token"`
		Refresh string `json:"refresh_token"`
	}

	authCodeRedirectUri := fmt.Sprintf("%s://%s/api/token", a.clientScheme, a.clientHostname)

	data := url.Values{}
	data.Add("grant_type", "authorization_code")
	data.Add("client_id", a.cognitoClientId)
	data.Add("redirect_uri", authCodeRedirectUri)
	data.Add("code", authCode)

	req, err := http.NewRequest("POST", "https://auth.expense.mleone.dev/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return domain.AuthTokens{}, fmt.Errorf("failed to construct token request: %w", err)
	}

	basicAuth := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", a.cognitoClientId, a.cognitoClientSecret))))
	req.Header.Add("Authorization", basicAuth)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := a.client.Do(req)
	if err != nil {
		return domain.AuthTokens{}, fmt.Errorf("request to token endpoint failed: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		return domain.AuthTokens{}, fmt.Errorf("received non-OK response from token endpoint: %v (%v)", res.Status, string(body))
	}

	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return domain.AuthTokens{}, fmt.Errorf("failed to read token response body: %w", err)
	}

	var tokens tokenResponse
	json.Unmarshal(bodyBytes, &tokens)

	return domain.AuthTokens{
		AccessToken: tokens.Access,
	}, nil
}

// ValidateToken confirms that the given token is correctly formatted and still valid.  A parsed, validated token is returned if so.
func (a *AuthClient) TokenIsValid(rawToken string) (jwt.Token, error) {
	const jwkUrl = "https://cognito-idp.us-east-1.amazonaws.com/us-east-1_hQXVbBbyZ/.well-known/jwks.json"

	keySet, err := jwk.Fetch(context.Background(), jwkUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch json web keys: %w", err)
	}

	token, err := jwt.Parse([]byte(rawToken), jwt.WithKeySet(keySet), jwt.WithValidate(true))
	if err != nil {
		return nil, fmt.Errorf("failed to parse provided token: %w", err)
	}

	return token, nil
}

// UserInfo is a container for various user data extrapolated from a given OAuth Token.
type UserInfo struct {
	Name       string
	ProfileUrl string
}

// GetUserInfo uses the provided token to request user info from Cognito.
func (a *AuthClient) GetUserInfo(authToken string) (UserInfo, error) {
	type userInfoResponse struct {
		Name       string `json:"name"`
		ProfileUrl string `json:"picture"`
	}

	req, err := http.NewRequest("GET", "https://auth.expense.mleone.dev/oauth2/userInfo", nil)
	if err != nil {
		return UserInfo{}, fmt.Errorf("failed to construct user info request: %w", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authToken))

	res, err := a.client.Do(req)
	if err != nil {
		return UserInfo{}, fmt.Errorf("request to user info endpoint failed: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		return UserInfo{}, fmt.Errorf("received non-OK response from user info endpoint: %v (%v)", res.Status, string(body))
	}

	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return UserInfo{}, fmt.Errorf("failed to read user info response body: %w", err)
	}

	var userInfo userInfoResponse
	json.Unmarshal(bodyBytes, &userInfo)

	return UserInfo(userInfo), nil
}
