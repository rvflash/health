package health

import (
	"context"
	"time"
)

// DefaultTimeout is the default maximum duration to check something before timeout.
const DefaultTimeout = 10 * time.Second

// Check is the interface of probe's check.
type Check func(ctx context.Context) error

// Probe represents a probe.
type Probe struct {
	Strategy Strategy
	Timeout  time.Duration
	Name     string
	Check    Check
}

// Strategy is a probe strategy.
type Strategy uint8

const (
	// Liveness is a strategy that indicates if the check fails that this instance is unhealthy and
	// should be destroyed or restarted.
	Liveness Strategy = iota
	// Readiness is a strategy that indicates if the probe fails that this application
	// should no longer receive any traffic.
	Readiness
)

// Error implements the error interface.
func (s Strategy) Error() string {
	switch s {
	case Readiness:
		return "health: readiness probe"
	}
	return "health: liveness probe"
}
