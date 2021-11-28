package api

import (
	"encoding/json"
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
	type response struct {
		Token string `json:"token"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := json.Marshal(response{Token: "tokenValue"})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(res)
	}
}
