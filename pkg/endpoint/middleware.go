package endpoint

import (
	"context"
	"fmt"
	"time"

	endpoint "github.com/go-kit/kit/endpoint"
	log "github.com/go-kit/kit/log"
	metrics "github.com/go-kit/kit/metrics"
)

// InstrumentingMiddleware returns an endpoint middleware that records
// the duration of each invocation to the passed histogram. The middleware adds
// a single field: "success", which is "true" if no error is returned, and
// "false" otherwise.
func InstrumentingMiddleware(duration metrics.Histogram) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			fmt.Println("------InstrumentingMiddleware[endpoint.Middleware]----------")

			defer func(begin time.Time) {
				fmt.Printf("Seconds: %v\n", time.Since(begin).Seconds())

				duration.With("success", fmt.Sprint(err == nil)).Observe(time.Since(begin).Seconds())

				fmt.Printf("duration: %v\n", duration)

			}(time.Now())
			return next(ctx, request)
		}
	}
}

// LoggingMiddleware returns an endpoint middleware that logs the
// duration of each invocation, and the resulting error, if any.
func LoggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			fmt.Println("------LoggingMiddleware[endpoint.Middleware]----------")
			defer func(begin time.Time) {
				logger.Log("transport_error", err, "took", time.Since(begin))
			}(time.Now())
			return next(ctx, request)
		}
	}
}

// HiMiddleware returns an endpoint middleware
func HiMiddleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			// Add your middleware logic here
			fmt.Println("------HiMiddleware----------")
			return next(ctx, request)
		}
	}
}
