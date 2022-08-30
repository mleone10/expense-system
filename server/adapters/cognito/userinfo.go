package cognito

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/mleone10/expense-system/domain"
)

// GetUserInfo uses the provided token to request user info from Cognito.
func (a *CognitoClient) GetUserInfo(ctx context.Context, authToken string) (domain.User, error) {
	type userInfoResponse struct {
		Name       string `json:"name"`
		ProfileUrl string `json:"picture"`
	}

	req, err := http.NewRequestWithContext(ctx, "GET", "https://auth.expense.mleone.dev/oauth2/userInfo", nil)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to construct user info request: %w", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authToken))

	res, err := a.client.Do(req)
	if err != nil {
		return domain.User{}, fmt.Errorf("request to user info endpoint failed: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		return domain.User{}, fmt.Errorf("received non-OK response from user info endpoint: %v (%v)", res.Status, string(body))
	}

	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to read user info response body: %w", err)
	}

	var userInfo userInfoResponse
	json.Unmarshal(bodyBytes, &userInfo)

	return domain.User{
		Name:       userInfo.Name,
		ProfileUrl: userInfo.ProfileUrl,
	}, nil
}
