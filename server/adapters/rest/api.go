package rest

// const urlParamOrgId string = "orgId"

// func NewServer(c Config) (Server, error) {

// 	s.router.Route("/api", func(r chi.Router) {

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
