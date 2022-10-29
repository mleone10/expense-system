package rest

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mleone10/expense-system/domain"
)

const urlParamOrgId string = "orgId"

func (hs *HttpServer) orgsRouter() func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/", hs.handleGetOrgs())
		r.Post("/", hs.handleCreateNewOrg())
		r.Route(fmt.Sprintf("/{%s}", urlParamOrgId), func(r chi.Router) {
			r.Get("/", hs.handleGetOrg())
			r.Patch("/", hs.handleUpdateOrg())
			r.Delete("/", hs.handleDeleteOrg())
		})
	}
}

func (hs HttpServer) handleGetOrgs() http.HandlerFunc {
	type org struct {
		Name  string `json:"name"`
		Id    string `json:"id"`
		Admin bool   `json:"admin"`
	}

	type response struct {
		Orgs []org `json:"orgs"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := getUserId(r)

		orgs, err := hs.orgService.GetOrgsForUser(r.Context(), userId)
		if err != nil {
			hs.writeError(w, r, err)
		}

		res := response{Orgs: []org{}}
		for _, o := range orgs {
			res.Orgs = append(res.Orgs, org{
				Name:  o.Name,
				Id:    string(o.Id),
				Admin: o.IsAdmin(userId),
			})
		}

		writeResponse(w, res)
	})
}

func (hs HttpServer) handleCreateNewOrg() http.HandlerFunc {
	type request struct {
		Name string `json:"name"`
	}

	type response struct {
		Id string `json:"id"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := readRequest(r, &req); err != nil {
			hs.writeError(w, r, fmt.Errorf("failed to read org creation request: %w", err))
			return
		}

		userId := getUserId(r)

		org, err := hs.orgService.CreateOrg(r.Context(), req.Name, userId)
		if err == domain.ErrInvalidRequest {
			hs.writeClientError(w, r, fmt.Errorf("failed to create org: %w", err))
			return
		}
		if err != nil {
			hs.writeError(w, r, fmt.Errorf("failed to create org: %w", err))
			return
		}

		writeResponse(w, response{Id: string(org.Id)})
	})
}

func (hs *HttpServer) handleGetOrg() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := getUserId(r)
		orgId := chi.URLParam(r, urlParamOrgId)

		org, err := hs.orgService.GetOrg(r.Context(), userId, domain.OrgId(orgId))
		if err != nil {
			hs.writeError(w, r, fmt.Errorf("failed to get org: %w", err))
		}

		writeResponse(w, org)
	})
}

func (hs *HttpServer) handleUpdateOrg() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}
func (hs *HttpServer) handleDeleteOrg() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}
