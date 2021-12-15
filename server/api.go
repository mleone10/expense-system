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

type ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request) error

func NewServer(c Config) (Server, error) {
	authClient, err := NewAuthClient(c)
	if err != nil {
		return Server{}, fmt.Errorf("failed to initialize auth client: %w", err)
	}

	s := Server{
		auth:   authClient,
		logger: *log.New(os.Stderr, "", log.LstdFlags),
	}

	s.mux.Handle("/api/token", s.commonMiddleware(s.handleToken()))

	return s, nil
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s Server) handleToken() ErrorHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ats, err := s.auth.GetAuthTokens(r.URL.Query().Get("code"))
		if err != nil {
			return fmt.Errorf("Failed to get auth tokens: %w", err)
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "authToken",
			Value:    ats.accessToken,
			HttpOnly: true,
			Expires:  time.Now().Add(time.Hour * 168),
		})

		return nil
	}
}

func (s Server) commonMiddleware(f ErrorHandlerFunc) http.HandlerFunc {
	return s.logRequests(s.errorHandler(f))
}

// TODO: Extend this to handle writing response objects to http.ResponseWriter and set default status code.
// Might need to extend http.ResponseWriter to support isWritten flag, etc.
func (s Server) errorHandler(f ErrorHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			s.logger.Printf("Request <> %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
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
