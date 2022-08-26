package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mleone10/expense-system/domain"
)

type keyTypeUserId string

const keyUserId keyTypeUserId = "userId"

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

		req := r.WithContext(context.WithValue(r.Context(), keyUserId, validatedToken.Subject()))
		next.ServeHTTP(w, req)
	})
}

func getUserId(r *http.Request) domain.UserId {
	userId := r.Context().Value(keyUserId)
	if userId != nil {
		return userId.(domain.UserId)
	}
	return ""
}
