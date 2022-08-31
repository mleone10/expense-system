package rest

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (hs *HttpServer) userRouter() func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/", hs.handleGetUser())
	}
}

func (hs HttpServer) handleGetUser() http.HandlerFunc {
	type response struct {
		Name       string `json:"name"`
		ProfileUrl string `json:"profileUrl"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken := getAuthToken(r)

		userInfo, err := hs.authenticatedUserService.GetAuthenticatedUserInfo(r.Context(), authToken)
		if err != nil {
			hs.writeError(w, r, fmt.Errorf("failed to get user info from identity provider: %w", err))
			return
		}

		writeResponse(w, response{Name: userInfo.Name, ProfileUrl: userInfo.ProfileUrl})
	})
}
