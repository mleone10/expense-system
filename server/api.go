package api

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Server struct {
	auth *authClient
	mux  http.ServeMux
}

func NewServer() (Server, error) {
	authClient, err := NewAuthClient()
	if err != nil {
		return Server{}, fmt.Errorf("failed to initialize auth client: %w", err)
	}

	s := Server{
		auth: authClient,
	}

	s.mux.Handle("/api/token", s.handleToken())

	return s, nil
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: Extract request logging to a middleware.
	log.Println("Handling request for", r.Method, r.RequestURI)
	s.mux.ServeHTTP(w, r)
}

func (s Server) handleToken() http.HandlerFunc {
	type response struct {
		Token string `json:"token"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ats, err := s.auth.GetAuthTokens(r.URL.Query().Get("code"))
		if err != nil {
			// TODO: Create common error logger/response writer.
			log.Printf("Failed to get auth tokens: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "authToken",
			Value:    ats.accessToken,
			HttpOnly: true,
			Expires:  time.Now().Add(time.Hour * 168),
		})
	}
}
