package rest

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"os"
// 	"time"

// 	"github.com/go-chi/chi/v5"
// 	"github.com/gofrs/uuid"
// )

// type Server struct {
// 	auth   *authClient
// 	router chi.Router
// 	logger *log.Logger
// 	orgs   orgRepo
// }

// type keyTypeRequestId string

// const keyRequestId keyTypeRequestId = "requestId"

// const urlParamOrgId string = "orgId"

// const testAdminUserId string = "nonProdTestAdmin"

// func NewServer(c Config) (Server, error) {

// 	tokenVerifierMiddleware := s.verifyToken
// 	if c.getSkipAuth() {
// 		tokenVerifierMiddleware = s.noOpTokenVerifier
// 	}

// 	s.router.Route("/api", func(r chi.Router) {
// 		r.Use(s.requestId)
// 		r.Use(s.logRequests)

// 		r.Group(func(r chi.Router) {
// 			r.Route("/orgs", func(r chi.Router) {
// 				r.Post("/", s.handleCreateNewOrg())
// 				r.Route(fmt.Sprintf("/{%s}", urlParamOrgId), func(r chi.Router) {
// 				})
// 			})

// 			r.Route("/user", func(r chi.Router) {
// 				r.Get("/", s.handleGetUser())
// 			})
// 		})
// 	})

// 	return s, nil
// }

// func (s Server) handleCreateNewOrg() http.HandlerFunc {
// 	type request struct {
// 		Name string `json:"name"`
// 	}

// 	type response struct {
// 		Id string `json:"id"`
// 	}

// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		var req request
// 		if err := s.readRequest(r, &req); err != nil {
// 			s.error(w, r, fmt.Errorf("failed to read org creation request: %w", err))
// 			return
// 		}

// 		userId, err := s.getUserId(r)
// 		if err != nil {
// 			s.error(w, r, fmt.Errorf("failed to get user id from request: %w", err))
// 			return
// 		}

// 		id, err := s.orgs.createOrg(req.Name, userId)
// 		if err != nil {
// 			s.error(w, r, fmt.Errorf("failed to create org with name %v: %w", req.Name, err))
// 			return
// 		}

// 		s.writeResponse(w, response{Id: id})
// 	})
// }

// func (s Server) handleGetUser() http.HandlerFunc {
// 	type response struct {
// 		Name       string `json:"name"`
// 		ProfileUrl string `json:"profileUrl"`
// 	}

// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		authCookie, err := r.Cookie(cookieNameAuthToken)
// 		if err != nil {
// 			s.error(w, r, fmt.Errorf("failed to get auth cookie from request: %w", err))
// 			return
// 		}

// 		userInfo, err := s.auth.GetUserInfo(authCookie.Value)
// 		if err != nil {
// 			s.error(w, r, fmt.Errorf("failed to get user info from identity provider: %w", err))
// 			return
// 		}

// 		s.writeResponse(w, response(userInfo))
// 	})
// }

// func (s Server) logRequests(next http.Handler) http.Handler {
// 	startTime := func() time.Time {
// 		return time.Now()
// 	}
// 	logReturn := func(startTime time.Time, r *http.Request) {
// 		s.logger.Printf("Request: %s completed in %s", s.getRequestId(r), time.Since(startTime))
// 	}

// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		defer logReturn(startTime(), r)
// 		s.logger.Printf("Request: %s Method: %s URI: %s", s.getRequestId(r), r.Method, r.RequestURI)
// 		next.ServeHTTP(w, r)
// 	})
// }

// func (s Server) requestId(next http.Handler) http.Handler {
// 	genRequestId := func() string {
// 		requestUuid, err := uuid.NewV4()
// 		if err != nil {
// 			s.logger.Println("Failed to generate request ID: %w", err)
// 			requestUuid = uuid.FromStringOrNil("")
// 		}
// 		return requestUuid.String()
// 	}

// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		req := r.WithContext(context.WithValue(r.Context(), keyRequestId, genRequestId()))
// 		next.ServeHTTP(w, req)
// 	})
// }

// func (s Server) getRequestId(r *http.Request) string {
// 	requestId := r.Context().Value(keyRequestId)
// 	if id, ok := requestId.(string); ok {
// 		return id
// 	} else {
// 		return "<no request id found>"
// 	}
// }

// func (s Server) noOpTokenVerifier(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		req := r.WithContext(context.WithValue(r.Context(), keyUserId, testAdminUserId))
// 		next.ServeHTTP(w, req)
// 	})
// }
