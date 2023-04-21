package health_test

import (
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/rvflash/health"
)

func TestStrategy_Error(t *testing.T) {
	t.Parallel()

	t.Run("Liveness", func(t *testing.T) {
		t.Parallel()
		is.New(t).True(strings.Contains(health.Liveness.Error(), "liveness"))
	})

	t.Run("Readiness", func(t *testing.T) {
		t.Parallel()
		is.New(t).True(strings.Contains(health.Readiness.Error(), "readiness"))
	})
}
