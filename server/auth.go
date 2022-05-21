package api

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
)

type authClient struct {
	client          http.Client
	cognitoClientId string
	basicAuth       string
	redirectUri     string
}

type authTokens struct {
	idToken      string
	accessToken  string
	refreshToken string
}

type authClientConfig interface {
	getClientHostname() string
	getClientScheme() string
	getCognitoClientId() string
	getCognitoClientSecret() string
}

type accessToken struct {
	validatedToken jwt.Token
}

const jwkUrl string = "https://cognito-idp.us-east-1.amazonaws.com/us-east-1_hQXVbBbyZ/.well-known/jwks.json"

func NewAuthClient(c authClientConfig) (*authClient, error) {

	a := authClient{
		client:          http.Client{},
		cognitoClientId: c.getCognitoClientId(),
		basicAuth:       fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.getCognitoClientId(), c.getCognitoClientSecret())))),
		redirectUri:     fmt.Sprintf("%s://%s/auth/callback", c.getClientScheme(), c.getClientHostname()),
	}

	return &a, nil
}

func (a *authClient) GetAuthTokens(authCode string) (authTokens, error) {
	type tokenResponse struct {
		Id      string `json:"id_token"`
		Access  string `json:"access_token"`
		Refresh string `json:"refresh_token"`
	}

	data := url.Values{}
	data.Add("grant_type", "authorization_code")
	data.Add("client_id", a.cognitoClientId)
	data.Add("redirect_uri", a.redirectUri)
	data.Add("code", authCode)

	req, err := http.NewRequest("POST", "https://auth.expense.mleone.dev/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return authTokens{}, fmt.Errorf("failed to construct token request: %w", err)
	}

	req.Header.Add("Authorization", a.basicAuth)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := a.client.Do(req)
	if err != nil {
		return authTokens{}, fmt.Errorf("request to token endpoint failed: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		return authTokens{}, fmt.Errorf("received non-OK response from token endpoint: %v (%v)", res.Status, string(body))
	}

	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return authTokens{}, fmt.Errorf("failed to read token response body: %w", err)
	}

	var tokens tokenResponse
	json.Unmarshal(bodyBytes, &tokens)

	return authTokens{accessToken: tokens.Access}, nil
}

func (a *authClient) TokenIsValid(rawToken string) (accessToken, error) {
	keySet, err := jwk.Fetch(context.Background(), jwkUrl)
	if err != nil {
		return accessToken{}, fmt.Errorf("failed to fetch json web keys: %w", err)
	}

	token, err := jwt.Parse([]byte(rawToken), jwt.WithKeySet(keySet))
	if err != nil {
		return accessToken{}, fmt.Errorf("failed to parse provided token: %w", err)
	}

	return accessToken{validatedToken: token}, nil
}

func (t accessToken) UserId() string {
	return t.validatedToken.Subject()
}

type UserInfo struct {
	Name       string
	ProfileUrl string
}

func (a *authClient) GetUserInfo(authToken string) (UserInfo, error) {
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

	return UserInfo{Name: userInfo.Name, ProfileUrl: userInfo.ProfileUrl}, nil
}
