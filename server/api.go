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
	logger *log.Logger
	orgs   orgRepo
}

type keyTypeRequestId string
type keyTypeUserId string

const keyRequestId keyTypeRequestId = "requestId"
const keyUserId keyTypeUserId = "userId"

const cookieNameAuthToken string = "authToken"

const urlParamOrgId string = "orgId"

const testAdminUserId string = "nonProdTestAdmin"

func NewServer(c Config) (Server, error) {
	authClient, err := NewAuthClient(c)
	if err != nil {
		return Server{}, fmt.Errorf("failed to initialize auth client: %w", err)
	}

	orgRepo, err := NewOrgRepo()
	if err != nil {
		return Server{}, fmt.Errorf("failed to initialize org repo: %w", err)
	}

	s := Server{
		auth:   authClient,
		router: chi.NewRouter(),
		logger: log.New(os.Stderr, "", log.LstdFlags),
		orgs:   orgRepo,
	}

	tokenVerifierMiddleware := s.verifyToken
	if c.getSkipAuth() {
		tokenVerifierMiddleware = s.noOpTokenVerifier
	}

	s.router.Route("/api", func(r chi.Router) {
		r.Use(s.requestId)
		r.Use(s.logRequests)

		r.Get("/health", s.handleHealth())
		r.Get("/token", s.handleToken())
		r.Get("/sign-out", s.handleSignOut())

		r.Group(func(r chi.Router) {
			r.Use(tokenVerifierMiddleware)

			r.Route("/orgs", func(r chi.Router) {
				r.Get("/", s.handleGetOrgs())
				r.Post("/", s.handleCreateNewOrg())
				r.Route(fmt.Sprintf("/{%s}", urlParamOrgId), func(r chi.Router) {
					r.Get("/", s.handleGetOrg())
					r.Post("/", s.handleUpdateOrg())
					r.Delete("/", s.handleDeleteOrg())
				})
			})

			r.Route("/user", func(r chi.Router) {
				r.Get("/", s.handleGetUser())
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
		s.writeResponse(w, response{
			Status: "ok",
		})
	})
}

func (s Server) handleToken() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ats, err := s.auth.GetAuthTokens(r.URL.Query().Get("code"))
		if err != nil {
			s.error(w, r, fmt.Errorf("failed to get auth tokens: %w", err))
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     cookieNameAuthToken,
			Value:    ats.accessToken,
			HttpOnly: true,
			Expires:  time.Now().Add(time.Hour * 168),
		})
		http.Redirect(w, r, s.auth.RedirectUrl(), http.StatusMovedPermanently)
	})
}

func (s Server) handleSignOut() http.HandlerFunc {
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

func (s Server) handleGetOrgs() http.HandlerFunc {
	type org struct {
		Name  string `json:"name"`
		Id    string `json:"id"`
		Admin bool   `json:"admin"`
	}

	type response struct {
		Orgs []org `json:"orgs"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, err := s.getUserId(r)
		if err != nil {
			s.error(w, r, fmt.Errorf("failed to get user id from request: %w", err))
			return
		}

		orgs, err := s.orgs.getOrgsForUser(userId)
		if err != nil {
			s.error(w, r, fmt.Errorf("failed to retrieve orgs for user %v: %w", userId, err))
			return
		}

		res := response{Orgs: []org{}}
		for _, o := range orgs {
			res.Orgs = append(res.Orgs, org{
				Name:  o.Name,
				Id:    o.Id,
				Admin: o.IsAdmin(),
			})
		}

		s.writeResponse(w, res)
	})
}

func (s Server) handleCreateNewOrg() http.HandlerFunc {
	type request struct {
		Name string `json:"name"`
	}

	type response struct {
		Id string `json:"id"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := s.readRequest(r, &req); err != nil {
			s.error(w, r, fmt.Errorf("failed to read org creation request: %w", err))
			return
		}

		userId, err := s.getUserId(r)
		if err != nil {
			s.error(w, r, fmt.Errorf("failed to get user id from request: %w", err))
			return
		}

		id, err := s.orgs.createOrg(req.Name, userId)
		if err != nil {
			s.error(w, r, fmt.Errorf("failed to create org with name %v: %w", req.Name, err))
			return
		}

		s.writeResponse(w, response{Id: id})
	})
}

func (s Server) handleGetOrg() http.HandlerFunc {
	type response struct {
		Id           string   `json:"id"`
		Name         string   `json:"name"`
		CreationDate string   `json:"creationDate"`
		Members      []string `json:"members"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := s.getUserId(r)
		if err != nil {
			s.error(w, r, fmt.Errorf("failed to get user id from request: %w", err))
			return
		}

		s.writeResponse(w, response{Members: []string{}})
	})
}

func (s Server) handleUpdateOrg() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}

func (s Server) handleDeleteOrg() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}

func (s Server) handleGetUser() http.HandlerFunc {
	type response struct {
		Name       string `json:"name"`
		ProfileUrl string `json:"profileUrl"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authCookie, err := r.Cookie(cookieNameAuthToken)
		if err != nil {
			s.error(w, r, fmt.Errorf("failed to get auth cookie from request: %w", err))
			return
		}

		userInfo, err := s.auth.GetUserInfo(authCookie.Value)
		if err != nil {
			s.error(w, r, fmt.Errorf("failed to get user info from identity provider: %w", err))
			return
		}

		s.writeResponse(w, response(userInfo))
	})
}

func (s Server) error(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, "internal server error", http.StatusInternalServerError)
	s.logger.Printf("Request: %s %v", s.getRequestId(r), err)
}

func (s Server) readRequest(r *http.Request, dest interface{}) error {
	return json.NewDecoder(r.Body).Decode(dest)
}

func (s Server) writeResponse(w http.ResponseWriter, src interface{}) error {
	return json.NewEncoder(w).Encode(src)
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
		req := r.WithContext(context.WithValue(r.Context(), keyRequestId, genRequestId()))
		next.ServeHTTP(w, req)
	})
}

func (s Server) getRequestId(r *http.Request) string {
	requestId := r.Context().Value(keyRequestId)
	if id, ok := requestId.(string); ok {
		return id
	} else {
		return "<no request id found>"
	}
}

func (s Server) getUserId(r *http.Request) (string, error) {
	userId := r.Context().Value(keyUserId)
	if id, ok := userId.(string); !ok {
		return "", fmt.Errorf("failed to convert user id to string")
	} else {
		return id, nil
	}
}

func (s Server) verifyToken(next http.Handler) http.Handler {
	markUnauthorized := func(w http.ResponseWriter) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie(cookieNameAuthToken)
		if err != nil {
			markUnauthorized(w)
			return
		}

		validatedToken, err := s.auth.TokenIsValid(tokenCookie.Value)
		if err != nil {
			s.error(w, r, fmt.Errorf("error while verifying auth token: %w", err))
			return
		}

		req := r.WithContext(context.WithValue(r.Context(), keyUserId, validatedToken.UserId()))
		next.ServeHTTP(w, req)
	})
}

func (s Server) noOpTokenVerifier(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := r.WithContext(context.WithValue(r.Context(), keyUserId, testAdminUserId))
		next.ServeHTTP(w, req)
	})
}
