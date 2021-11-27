package api

import (
	"net/http"
)

type Server struct {
	mux http.ServeMux
}

func NewServer() (Server, error) {
	s := Server{}
	s.mux.Handle("/token", s.handleToken())
	return s, nil
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s Server) handleToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:  "token",
			Value: "tokenValue",
		})
	}
}