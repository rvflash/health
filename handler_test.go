package health_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/rvflash/health"
)

const (
	name01 = "name_01"
	name02 = "name_02"
)

func TestHandlerFunc(t *testing.T) {
	t.Parallel()

	var (
		are = is.New(t)
		err = errors.New("oops")
	)
	for name, tc := range map[string]struct {
		// inputs
		in health.Checker
		// outputs
		code int
		out  []string
	}{
		"Default": {code: http.StatusOK},
		"Check OK": {
			in: []health.Probe{
				{
					Name: name01,
					Check: func(context.Context) error {
						return nil
					},
				},
			},
			code: http.StatusOK,
			out:  []string{`"status":"OK"`},
		},
		"Liveness failed": {
			in: []health.Probe{
				{
					Name: name01,
					Check: func(context.Context) error {
						return err
					},
				},
			},
			code: http.StatusServiceUnavailable,
			out:  []string{`"status":"Service Unavailable"`},
		},
		"Liveness timeout": {
			in: []health.Probe{
				{
					Name: name01,
					Check: func(context.Context) error {
						return context.DeadlineExceeded
					},
				},
			},
			code: http.StatusGatewayTimeout,
			out:  []string{`"status":"Gateway Timeout"`},
		},
		"Readiness failed": {
			in: []health.Probe{
				{
					Strategy: health.Readiness,
					Name:     name01,
					Check: func(context.Context) error {
						return err
					},
				},
			},
			code: http.StatusFailedDependency,
			out:  []string{`"status":"Failed Dependency"`},
		},
		"Readiness timeout": {
			in: []health.Probe{
				{
					Strategy: health.Readiness,
					Name:     name01,
					Check: func(context.Context) error {
						return context.DeadlineExceeded
					},
				},
			},
			code: http.StatusRequestTimeout,
			out:  []string{`"status":"Request Timeout"`},
		},
	} {
		tt := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req, err := http.NewRequest("GET", "/health", nil)
			are.NoErr(err)
			rec := httptest.NewRecorder()
			health.HandlerFunc(tt.in).ServeHTTP(rec, req)

			are.Equal(tt.code, rec.Code) // mismatch response code
			out := rec.Body.String()
			for _, s := range tt.out {
				are.True(strings.Contains(out, s)) // mismatch response
			}
		})
	}
}
