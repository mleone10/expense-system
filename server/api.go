package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gofrs/uuid"
)

type Server struct {
	auth   *authClient
	mux    http.ServeMux
	logger log.Logger
}

func NewServer(c Config) (Server, error) {
	authClient, err := NewAuthClient(c)
	if err != nil {
		return Server{}, fmt.Errorf("failed to initialize auth client: %w", err)
	}

	s := Server{
		auth:   authClient,
		logger: *log.New(os.Stderr, "", log.LstdFlags),
	}

	s.mux.Handle("/api/token", s.logRequests(s.handleToken()))

	return s, nil
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
			s.logger.Printf("Failed to get auth tokens: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "authToken",
			Value:    ats.accessToken,
			HttpOnly: true,
			Expires:  time.Now().Add(time.Hour * 168),
		})
		w.WriteHeader(http.StatusOK)
	}
}

func (s Server) logRequests(next http.HandlerFunc) http.HandlerFunc {
	startTime := func() time.Time {
		return time.Now()
	}
	logReturn := func(startTime time.Time, requestId string) {
		s.logger.Printf("Request: %s completed in %s", requestId, time.Since(startTime))
	}

	return func(w http.ResponseWriter, r *http.Request) {
		requestUuid, err := uuid.NewV4()
		if err != nil {
			s.logger.Println("Failed to generate request ID: %w", err)
			requestUuid = uuid.FromStringOrNil("")
		}

		defer logReturn(startTime(), requestUuid.String())
		s.logger.Printf("Request: %s Method: %s URI: %s", requestUuid.String(), r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	}
}
