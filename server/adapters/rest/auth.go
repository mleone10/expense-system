package rest

import (
	"fmt"
	"net/http"
	"time"
)

const cookieNameAuthToken string = "authToken"

func (hs HttpServer) handleToken() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ats, err := hs.authClient.GetAuthTokens(r.URL.Query().Get("code"))
		if err != nil {
			hs.writeError(w, r, fmt.Errorf("failed to get auth tokens: %w", err))
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     cookieNameAuthToken,
			Value:    ats.AccessToken,
			HttpOnly: true,
			Expires:  time.Now().Add(time.Hour * 168),
		})
		http.Redirect(w, r, hs.authClient.RedirectUrl(), http.StatusMovedPermanently)
	})
}

func (hs HttpServer) handleSignOut() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     cookieNameAuthToken,
			Value:    "",
			HttpOnly: true,
			Expires:  time.Now().Add(time.Hour * -1),
		})
		w.WriteHeader(http.StatusOK)
	})
}
