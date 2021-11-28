package api

import (
	"log"
	"net/http"
)

type Server struct {
	mux http.ServeMux
}

func NewServer() (Server, error) {
	s := Server{}
	s.mux.Handle("/api/token", s.handleToken())
	return s, nil
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling request for", r.Method, r.RequestURI)
	s.mux.ServeHTTP(w, r)
}

func (s Server) handleToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:  "token",
			Value: "tokenValue",
		})
		w.WriteHeader(http.StatusNoContent)
	}
}
