package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mleone10/expense-system/domain"
)

type (
	keyTypeAuthToken string
	keyTypeUserId    string
)

const (
	keyAuthToken keyTypeAuthToken = "authToken"
	keyUserId    keyTypeUserId    = "userId"
)

const (
	testAdminUserId string = "nonProdTestAdmin"
	testAuthToken   string = "nonProdTestToken"
)

func (hs HttpServer) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie(cookieNameAuthToken)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		validatedToken, err := hs.authClient.TokenIsValid(tokenCookie.Value)
		if err != nil {
			hs.writeError(w, r, fmt.Errorf("error while verifying auth token: %w", err))
			return
		}

		reqWithUserId := r.WithContext(context.WithValue(r.Context(), keyUserId, validatedToken.Subject()))
		reqWithAuthToken := r.WithContext(context.WithValue(reqWithUserId.Context(), keyAuthToken, tokenCookie.Value))
		next.ServeHTTP(w, reqWithAuthToken)
	})
}

func getAuthToken(r *http.Request) string {
	authToken := r.Context().Value(keyAuthToken)
	if authToken != nil {
		return authToken.(string)
	}
	return ""
}

func getUserId(r *http.Request) domain.UserId {
	userId := r.Context().Value(keyUserId)
	if userId != nil {
		return domain.UserId(userId.(string))
	}
	return ""
}

func (hs HttpServer) noOpAuthVerifier(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqWithUserId := r.WithContext(context.WithValue(r.Context(), keyUserId, testAdminUserId))
		next.ServeHTTP(w, reqWithUserId)
	})
}
