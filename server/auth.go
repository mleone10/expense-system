package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
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

func NewAuthClient() (*authClient, error) {
	cognitoClientId := os.Getenv("COGNITO_CLIENT_ID")

	a := authClient{
		client:          http.Client{},
		cognitoClientId: cognitoClientId,
		basicAuth:       fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", cognitoClientId, os.Getenv("COGNITO_CLIENT_SECRET"))))),
		redirectUri:     fmt.Sprintf("%s://%s/auth/callback", os.Getenv("CLIENT_SCHEME"), os.Getenv("CLIENT_HOST")),
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
		return authTokens{}, fmt.Errorf("received non-OK response from token endpoint: %v", res.Status)
	}

	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var tokens tokenResponse
	json.Unmarshal(bodyBytes, &tokens)

	return authTokens{accessToken: tokens.Access}, nil
}
