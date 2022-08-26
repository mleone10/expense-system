package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
)

type keyTypeRequestId string

const keyRequestId keyTypeRequestId = "requestId"

func (hs HttpServer) requestIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uuid, err := uuid.NewV4()
		if err != nil {
			hs.writeError(w, r, fmt.Errorf("failed to generate request id: %w", err))
			return
		}

		req := r.WithContext(context.WithValue(r.Context(), keyRequestId, uuid.String()))
		next.ServeHTTP(w, req)
	})
}

func (hs HttpServer) getRequestId(r *http.Request) string {
	userId := r.Context().Value(keyRequestId)
	if userId != nil {
		return userId.(string)
	}
	return ""
}
