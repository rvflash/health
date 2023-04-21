// Package health provides healthcheck functionality to monitor application's components.
// An `net/http` handler is also exposed and provides a response in JSON format with details about the result.
//
//	{
//		"date": "2006-01-02T15:04:05Z07:00",
//		"latency": "23ms",
//		"status": "OK",
//		"errors": "probe: error message"
//	}
package health

import (
	"context"
	"fmt"
	"time"

	"github.com/rvflash/workr"
)

// Checker is a list of probe to launch to verify the application's integrity.
type Checker []Probe

// New creates a Checker based on a list of Probe.
func New(probes ...Probe) Checker {
	return probes
}

// Do launches in concurrency all the probes to check.
func (c Checker) Do(ctx context.Context) error {
	w := workr.New(workr.ReturnAllErrors())
	for _, probe := range c {
		p := probe
		w.Go(func() error {
			err := checkWithTimeout(ctx, p.Timeout, p.Check)
			if err != nil {
				return fmt.Errorf("%s: %w: %w", p.Name, err, p.Strategy)
			}
			return nil
		})
	}
	return w.Wait()
}

// We voluntarily do not wait for the response of the Check,
// in order to protect the process against probes not listening the context.
// So the go routine can leak.
func checkWithTimeout(parent context.Context, timeout time.Duration, f Check) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()
	c := make(chan error, 1)
	go func() { c <- f(ctx) }()
	select {
	case <-ctx.Done():
		// <-c leaks in order to return as soon as possible the problem.
		return ctx.Err()
	case err := <-c:
		return err
	}
}
