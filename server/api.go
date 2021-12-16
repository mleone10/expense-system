package api

import (
	"context"
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

type KeyTypeRequestId string

const KeyRequestId KeyTypeRequestId = "requestId"

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
	return s.injectRequestId(s.logRequests(s.errorHandler(f)))
}

// TODO: Extend this to handle writing response objects to http.ResponseWriter and set default status code.
// Might need to extend http.ResponseWriter to support isWritten flag, etc.
func (s Server) errorHandler(f ErrorHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			s.logger.Printf("Request: %s %v", s.getRequestId(r), err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (s Server) logRequests(next http.HandlerFunc) http.HandlerFunc {
	startTime := func() time.Time {
		return time.Now()
	}
	logReturn := func(startTime time.Time, r *http.Request) {
		s.logger.Printf("Request: %s completed in %s", s.getRequestId(r), time.Since(startTime))
	}

	return func(w http.ResponseWriter, r *http.Request) {
		defer logReturn(startTime(), r)
		s.logger.Printf("Request: %s Method: %s URI: %s", s.getRequestId(r), r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	}
}

func (s Server) injectRequestId(next http.HandlerFunc) http.HandlerFunc {
	genRequestId := func() string {
		requestUuid, err := uuid.NewV4()
		if err != nil {
			s.logger.Println("Failed to generate request ID: %w", err)
			requestUuid = uuid.FromStringOrNil("")
		}
		return requestUuid.String()
	}

	return func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), KeyRequestId, genRequestId())))
	}
}

func (s Server) getRequestId(r *http.Request) string {
	requestId := r.Context().Value(KeyRequestId)
	if id, ok := requestId.(string); ok {
		return id
	} else {
		s.logger.Println("No request ID found")
		return ""
	}
}
