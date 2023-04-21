package health

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

// HandlerFunc returns an HTTP HandlerFunc exposing the Checker's result as a JSON string.
func HandlerFunc(c Checker) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var (
			res = response{Date: time.Now()}
			err = c.Do(req.Context())
		)
		if err != nil {
			res.Errors = err.Error()
		}
		code := statusCode(err)
		res.Status = http.StatusText(code)
		res.Latency = time.Since(res.Date).String()
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(res)
	}
}

func statusCode(err error) int {
	if err != nil {
		if errors.Is(err, Liveness) {
			if errors.Is(err, context.DeadlineExceeded) {
				return http.StatusGatewayTimeout
			}
			return http.StatusServiceUnavailable
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return http.StatusRequestTimeout
		}
		return http.StatusFailedDependency
	}
	return http.StatusOK
}

type response struct {
	Date    time.Time `json:"date,omitempty"`
	Latency string    `json:"latency,omitempty"`
	Status  string    `json:"status"`
	Errors  string    `json:"errors,omitempty"`
}
