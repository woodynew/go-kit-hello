package thrift

import (
	"context"
	"errors"
	endpoint1 "hello/pkg/endpoint"
	service "hello/pkg/service"
	"time"

	"github.com/go-kit/kit/circuitbreaker"
	endpoint "github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"

	addthrift "hello/pkg/thrift/gen-go/addsvc"
)

// NewThriftClient returns an AddService backed by a Thrift server described by
// the provided client. The caller is responsible for constructing the client,
// and eventually closing the underlying transport. We bake-in certain middlewares,
// implementing the client library pattern.
func NewThriftClient(client *addthrift.AddServiceClient) service.HelloService {
	// We construct a single ratelimiter middleware, to limit the total outgoing
	// QPS from this client to all methods on the remote instance. We also
	// construct per-endpoint circuitbreaker middlewares to demonstrate how
	// that's done, although they could easily be combined into a single breaker
	// for the entire remote instance, too.
	limiter := ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))

	// Each individual endpoint is an http/transport.Client (which implements
	// endpoint.Endpoint) that gets wrapped with various middlewares. If you
	// could rely on a consistent set of client behavior.
	var sumEndpoint endpoint.Endpoint
	{
		sumEndpoint = MakeThriftSumEndpoint(client)
		sumEndpoint = limiter(sumEndpoint)
		sumEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "Sum",
			Timeout: 30 * time.Second,
		}))(sumEndpoint)
	}

	// The Concat endpoint is the same thing, with slightly different
	// middlewares to demonstrate how to specialize per-endpoint.
	var concatEndpoint endpoint.Endpoint
	{
		concatEndpoint = MakeThriftConcatEndpoint(client)
		concatEndpoint = limiter(concatEndpoint)
		concatEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "Concat",
			Timeout: 10 * time.Second,
		}))(concatEndpoint)
	}

	// The Concat endpoint is the same thing, with slightly different
	// middlewares to demonstrate how to specialize per-endpoint.
	var fooEndpoint endpoint.Endpoint
	{
		fooEndpoint = MakeThriftFooEndpoint(client)
		fooEndpoint = limiter(fooEndpoint)
		fooEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "Foo",
			Timeout: 10 * time.Second,
		}))(fooEndpoint)
	}

	// Returning the endpoint.Set as a service.Service relies on the
	// endpoint.Set implementing the Service methods. That's just a simple bit
	// of glue code.
	return endpoint1.Endpoints{
		FooEndpoint:    fooEndpoint,
		SumEndpoint:    sumEndpoint,
		ConcatEndpoint: concatEndpoint,
	}
}

// MakeThriftSumEndpoint returns an endpoint that invokes the passed Thrift client.
// Useful only in clients, and only until a proper transport/thrift.Client exists.
func MakeThriftFooEndpoint(client *addthrift.AddServiceClient) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(endpoint1.FooRequest)
		reply, err := client.Foo(ctx, string(req.S))
		if err == ErrIntOverflow {
			return nil, err // special case; see comment on ErrIntOverflow
		}
		return endpoint1.FooResponse{Rs: string(reply.Value), Err: err}, nil
	}
}

// MakeThriftSumEndpoint returns an endpoint that invokes the passed Thrift client.
// Useful only in clients, and only until a proper transport/thrift.Client exists.
func MakeThriftSumEndpoint(client *addthrift.AddServiceClient) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(endpoint1.SumRequest)
		reply, err := client.Sum(ctx, int64(req.A), int64(req.B))
		if err == ErrIntOverflow {
			return nil, err // special case; see comment on ErrIntOverflow
		}
		return endpoint1.SumResponse{I0: int(reply.Value), E1: err}, nil
	}
}

// MakeThriftConcatEndpoint returns an endpoint that invokes the passed Thrift
// client. Useful only in clients, and only until a proper
// transport/thrift.Client exists.
func MakeThriftConcatEndpoint(client *addthrift.AddServiceClient) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(endpoint1.ConcatRequest)
		reply, err := client.Concat(ctx, req.A, req.B)
		return endpoint1.ConcatResponse{S0: reply.Value, E1: err}, nil
	}
}

func str2err(s string) error {
	if s == "" {
		return nil
	}
	return errors.New(s)
}

func err2str(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

var (
	// ErrTwoZeroes is an arbitrary business rule for the Add method.
	ErrTwoZeroes = errors.New("can't sum two zeroes")

	// ErrIntOverflow protects the Add method. We've decided that this error
	// indicates a misbehaving service and should count against e.g. circuit
	// breakers. So, we return it directly in endpoints, to illustrate the
	// difference. In a real service, this probably wouldn't be the case.
	ErrIntOverflow = errors.New("integer overflow")

	// ErrMaxSizeExceeded protects the Concat method.
	ErrMaxSizeExceeded = errors.New("result exceeds maximum size")
)
