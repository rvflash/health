package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/rvflash/health"
)

func main() {
	c := health.New(
		health.Probe{
			Strategy: health.Liveness,
			Timeout:  time.Second,
			Name:     "file",
			Check:    health.CreateFileCheck("/tmp", "check"),
		},
		health.Probe{
			Strategy: health.Readiness,
			Timeout:  500 * time.Millisecond,
			Name:     "oops",
			Check: func(context.Context) error {
				time.Sleep(time.Second)
				return nil
			},
		},
	)
	http.HandleFunc("/health", health.HandlerFunc(c))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
