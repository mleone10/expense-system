package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
)

type Server struct {
	auth   *authClient
	router chi.Router
	logger log.Logger
}

type KeyTypeRequestId string

const KeyRequestId KeyTypeRequestId = "requestId"

func NewServer(c Config) (Server, error) {
	authClient, err := NewAuthClient(c)
	if err != nil {
		return Server{}, fmt.Errorf("failed to initialize auth client: %w", err)
	}

	s := Server{
		auth:   authClient,
		router: chi.NewRouter(),
		logger: *log.New(os.Stderr, "", log.LstdFlags),
	}

	s.router.Use(s.requestId)
	s.router.Use(s.logRequests)
	s.router.Route("/api", func(r chi.Router) {
		r.Get("/health", s.handleHealth())
		r.Get("/token", s.handleToken())
		r.Route("/orgs", handleOrgs())
		r.Route("/users", func(r chi.Router) {
			r.Route("/{userID}", func(r chi.Router) {
			})
		})
	})

	return s, nil
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s Server) handleHealth() http.HandlerFunc {
	type response struct {
		Status string `json:"status"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(response{
			Status: "ok",
		})
	})
}

func (s Server) handleToken() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ats, err := s.auth.GetAuthTokens(r.URL.Query().Get("code"))
		if err != nil {
			s.error(w, r, fmt.Errorf("failed to get auth tokens: %w", err))
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "authToken",
			Value:    ats.accessToken,
			HttpOnly: true,
			Expires:  time.Now().Add(time.Hour * 168),
		})
	})
}

func (s Server) error(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, "internal server error", http.StatusInternalServerError)
	s.logger.Printf("Request: %s %v", s.getRequestId(r), err)
}

func (s Server) logRequests(next http.Handler) http.Handler {
	startTime := func() time.Time {
		return time.Now()
	}
	logReturn := func(startTime time.Time, r *http.Request) {
		s.logger.Printf("Request: %s completed in %s", s.getRequestId(r), time.Since(startTime))
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer logReturn(startTime(), r)
		s.logger.Printf("Request: %s Method: %s URI: %s", s.getRequestId(r), r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func (s Server) requestId(next http.Handler) http.Handler {
	genRequestId := func() string {
		requestUuid, err := uuid.NewV4()
		if err != nil {
			s.logger.Println("Failed to generate request ID: %w", err)
			requestUuid = uuid.FromStringOrNil("")
		}
		return requestUuid.String()
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := r.WithContext(context.WithValue(r.Context(), KeyRequestId, genRequestId()))
		next.ServeHTTP(w, req)
	})
}

func (s Server) getRequestId(r *http.Request) string {
	requestId := r.Context().Value(KeyRequestId)
	if id, ok := requestId.(string); ok {
		return id
	} else {
		return "<no request id found>"
	}
}
