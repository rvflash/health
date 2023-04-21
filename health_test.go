package health_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/rvflash/health"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("Default", func(t *testing.T) {
		t.Parallel()
		is.New(t).Equal(nil, health.New())
	})

	t.Run("OK", func(t *testing.T) {
		t.Parallel()
		var (
			are = is.New(t)
			out = health.New(health.Probe{Name: name01}, health.Probe{Name: name02})
		)
		are.Equal(2, len(out))         // mismatch size
		are.Equal(name01, out[0].Name) // unexpected first probe
		are.Equal(name02, out[1].Name) // unexpected second probe
	})
}

func TestChecker_Do(t *testing.T) {
	t.Parallel()

	t.Run("Timeout exceeded", func(t *testing.T) {
		const max = 20 * time.Millisecond
		var (
			are = is.New(t)
			chk = health.New(health.Probe{
				Timeout: max,
				Check: func(ctx context.Context) error {
					time.Sleep(time.Second)
					return nil
				},
			})
			beg = time.Now()
			err = chk.Do(context.Background())
		)
		if time.Since(beg) > (max*20/100)+max {
			are.Fail() // Timeout too long
		}
		are.True(errors.Is(err, context.DeadlineExceeded))
	})
}
