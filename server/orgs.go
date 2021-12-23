package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const urlParamOrgId string = "orgId"

func handleOrgs() func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/", handleGetOrgs())
		r.Post("/", handleCreateNewOrg())
		r.Route(fmt.Sprintf("/{%s}", urlParamOrgId), func(r chi.Router) {
			r.Get("/", handleGetOrg())
			r.Post("/", handleUpdateOrg())
			r.Delete("/", handleDeleteOrg())
		})
	}
}

func handleGetOrgs() http.HandlerFunc {
	type response struct {
		Orgs []string `json:"orgs"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(response{
			Orgs: []string{},
		})
	})
}

func handleCreateNewOrg() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}

func handleGetOrg() http.HandlerFunc {
	type response struct {
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}

func handleUpdateOrg() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}

func handleDeleteOrg() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}
