# Health✓

[![GoDoc](https://godoc.org/github.com/rvflash/health?status.svg)](https://godoc.org/github.com/rvflash/health)
[![Build Status](https://github.com/rvflash/health/workflows/build/badge.svg)](https://github.com/rvflash/health/actions?workflow=build)
[![Code Coverage](https://codecov.io/gh/rvflash/health/branch/main/graph/badge.svg)](https://codecov.io/gh/rvflash/health)
[![Go Report Card](https://goreportcard.com/badge/github.com/rvflash/health?)](https://goreportcard.com/report/github.com/rvflash/health)


`health` is a Go package providing facilities to check liveness or readiness application dependencies.


## Features

1. Exposes an HTTP handler that retrieves health status of the application with HTTP status code and JSON response.
   ```json
   {
      "date":"2023-04-21T21:54:06.238997+02:00",
      "latency":"503.314958ms",
      "status":"Request Timeout",
      "errors":"mysql: context deadline exceeded: health: readiness probe"
   }
   ```
2. Provides a function interface to implement to check any dependency: `func(ctx context.Context) error`.
3. 2 strategies: 
   - `Liveness` to indicate if the probe failed that this instance is unhealthy and should be destroyed or restarted. <br />
   In case of error an HTTP status `ServiceUnavailable` is returned, `GatewayTimeout` in case of deadline exceeded.
   - `Readiness` to notify if the probe failed that this application should no longer receive any traffic. <br />
   In case of error an HTTP status `FailedDependency` is returned, `RequestTimeout` in case of deadline exceeded.
4. Each probe has a name, a timeout, a strategy and a function to check.
5. Provides a function to check file writing.


## Example

Here we create a health checker to verify MySQL database connection and NFS write access in directory named `/data`.

```go
   db, err := sql.Open("mysql", "user:password@/dbname")
   if err != nil {
      log.Fatal(err)
   }
   c := health.New(
      health.Probe{
         Strategy: health.Readiness,
         Timeout:  time.Second,
         Name:     "mysql",
         Check:    db.PingContext,
      },
      health.Probe{
         Strategy: health.Liveness,
         Timeout:  time.Second,
         Name:     "nfs",
         Check:    health.CreateFileCheck("/data", "check"),
      },
   )
   http.HandleFunc("/health", health.HandlerFunc(c))
   log.Fatal(http.ListenAndServe(":8080", nil))
```

> See example directory for another sample.