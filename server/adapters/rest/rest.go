package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/mleone10/expense-system/domain"
)

const (
	cookieNameAuthToken string = "authToken"
)

type HttpServer struct {
	router     chi.Router
	authClient domain.AuthClient
	logger     domain.Logger
}

type OptionFunc func(*HttpServer)

func NewServer(options ...OptionFunc) (*HttpServer, error) {
	hs := &HttpServer{
		router: chi.NewRouter(),
	}

	for _, opt := range options {
		opt(hs)
	}

	hs.router.Route("/api", func(r chi.Router) {
		// r.Use(s.requestId)
		// r.Use(s.logRequests)

		r.Get("/health", hs.handleHealth())
		r.Get("/token", hs.handleToken())
		r.Get("/sign-out", hs.handleSignOut())

		// r.Group(func(r chi.Router) {
		// 	r.Use(tokenVerifierMiddleware)

		// 	r.Route("/orgs", func(r chi.Router) {
		// 		r.Get("/", s.handleGetOrgs())
		// 		r.Post("/", s.handleCreateNewOrg())
		// 		r.Route(fmt.Sprintf("/{%s}", urlParamOrgId), func(r chi.Router) {
		// 			r.Get("/", s.handleGetOrg())
		// 			r.Post("/", s.handleUpdateOrg())
		// 			r.Delete("/", s.handleDeleteOrg())
		// 		})
		// 	})

		// 	r.Route("/user", func(r chi.Router) {
		// 		r.Get("/", s.handleGetUser())
		// 	})
		// })
	})

	return hs, nil
}

func WithAuthClient(authClient domain.AuthClient) OptionFunc {
	return func(hs *HttpServer) {
		hs.authClient = authClient
	}
}

func WithLogger(logger domain.Logger) OptionFunc {
	return func(hs *HttpServer) {
		hs.logger = logger
	}
}

func (hs HttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hs.router.ServeHTTP(w, r)
}

func (hs HttpServer) handleHealth() http.HandlerFunc {
	type response struct {
		Status string `json:"status"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeResponse(w, response{
			Status: "ok",
		})
	})
}

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

func readRequest(r *http.Request, dest interface{}) error {
	return json.NewDecoder(r.Body).Decode(dest)
}

func writeResponse(w http.ResponseWriter, src interface{}) error {
	return json.NewEncoder(w).Encode(src)
}

func (hs HttpServer) writeError(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, "internal server error", http.StatusInternalServerError)
	hs.logger.Print(context.Background(), err)
	// hs.logger.Printf("Request: %s %v", hs.getRequestId(r), err)
}
