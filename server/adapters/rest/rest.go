package rest

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/mleone10/expense-system/domain"
)

type HttpServer struct {
	router                   chi.Router
	authClient               domain.AuthClient
	orgService               domain.OrgService
	authenticatedUserService domain.AuthenticatedUserService
	logger                   domain.Logger
	activeAuthMiddleware     func(http.Handler) http.Handler
}

type OptionFunc func(*HttpServer)

func NewServer(options ...OptionFunc) (*HttpServer, error) {
	hs := &HttpServer{
		router: chi.NewRouter(),
	}
	hs.activeAuthMiddleware = hs.authMiddleware

	for _, opt := range options {
		opt(hs)
	}

	hs.router.Route("/api", func(r chi.Router) {
		r.Use(hs.requestIdMiddleware)
		r.Use(hs.requestLoggingMiddleware)

		r.Get("/health", hs.handleHealth())
		r.Get("/token", hs.handleToken())
		r.Get("/sign-out", hs.handleSignOut())

		r.Group(func(r chi.Router) {
			r.Use(hs.activeAuthMiddleware)

			r.Route("/orgs", hs.orgsRouter())
			r.Route("/user", hs.userRouter())
		})
	})

	return hs, nil
}

func WithAuthClient(authClient domain.AuthClient) OptionFunc {
	return func(hs *HttpServer) {
		hs.authClient = authClient
	}
}

func WithOrgService(orgService domain.OrgService) OptionFunc {
	return func(hs *HttpServer) {
		hs.orgService = orgService
	}
}

func WithAuthenticatedUserService(authenticatedUserService domain.AuthenticatedUserService) OptionFunc {
	return func(hs *HttpServer) {
		hs.authenticatedUserService = authenticatedUserService
	}
}

func WithLogger(logger domain.Logger) OptionFunc {
	return func(hs *HttpServer) {
		hs.logger = logger
	}
}

func WithSkipAuth() OptionFunc {
	return func(hs *HttpServer) {
		hs.activeAuthMiddleware = hs.noOpAuthVerifier
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

func readRequest(r *http.Request, dest interface{}) error {
	return json.NewDecoder(r.Body).Decode(dest)
}

func writeResponse(w http.ResponseWriter, src interface{}) error {
	return json.NewEncoder(w).Encode(src)
}

func (hs HttpServer) writeError(w http.ResponseWriter, r *http.Request, err error) {
	type errorPayload struct {
		RequestId string    `json:"requestId"`
		Time      time.Time `json:"time"`
		ErrorMsg  string    `json:"errorMsg"`
	}
	http.Error(w, "internal server error", http.StatusInternalServerError)
	hs.logger.Print(r.Context(), errorPayload{getRequestId(r), time.Now(), err.Error()})
}

func (hs HttpServer) writeClientError(w http.ResponseWriter, r *http.Request, err error) {
	type errorPayload struct {
		Error error `json:"error"`
	}
	http.Error(w, "invalid request", http.StatusBadRequest)
	hs.logger.Print(r.Context(), errorPayload{Error: err})
}
