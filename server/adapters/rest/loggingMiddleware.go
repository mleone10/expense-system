package rest

import (
	"net/http"
	"time"
)

type inboundLogPayload struct {
	RequestId string    `json:"requestId"`
	StartTime time.Time `json:"startTime"`
	Method    string    `json:"method"`
	Uri       string    `json:"uri"`
}

type outboundLogPayload struct {
	RequestId string        `json:"requestId"`
	EndTime   time.Time     `json:"endTime"`
	Duration  time.Duration `json:"duration"`
}

func (hs HttpServer) requestLoggingMiddleware(next http.Handler) http.Handler {
	logReturn := func(r *http.Request, startTime time.Time) {
		endTime := time.Now()

		hs.logger.Print(r.Context(), outboundLogPayload{
			RequestId: getRequestId(r),
			EndTime:   endTime,
			Duration:  endTime.Sub(startTime),
		})
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		payload := inboundLogPayload{
			RequestId: getRequestId(r),
			StartTime: startTime,
			Method:    r.Method,
			Uri:       r.RequestURI,
		}

		defer logReturn(r, startTime)
		hs.logger.Print(r.Context(), payload)
		next.ServeHTTP(w, r)
	})
}
