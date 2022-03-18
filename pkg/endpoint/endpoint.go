package endpoint

import (
	"context"

	service "github.com/woodynew/go-kit-hello/pkg/service"

	endpoint "github.com/go-kit/kit/endpoint"
)

// FooRequest collects the request parameters for the Foo method.
type FooRequest struct {
	S string `json:"s"`
}

// FooResponse collects the response parameters for the Foo method.
type FooResponse struct {
	Rs  string `json:"rs"`
	Err error  `json:"err"`
}

// MakeFooEndpoint returns an endpoint that invokes Foo on the service.
func MakeFooEndpoint(s service.HelloService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(FooRequest)
		rs, err := s.Foo(ctx, req.S)
		return FooResponse{
			Err: err,
			Rs:  rs,
		}, nil
	}
}

// Failed implements Failer.
func (r FooResponse) Failed() error {
	return r.Err
}

// Failure is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so they've
// failed, and if so encode them using a separate write path based on the error.
type Failure interface {
	Failed() error
}

// Foo implements Service. Primarily useful in a client.
func (e Endpoints) Foo(ctx context.Context, s string) (rs string, err error) {
	request := FooRequest{S: s}
	response, err := e.FooEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(FooResponse).Rs, response.(FooResponse).Err
}

// SumRequest collects the request parameters for the Sum method.
type SumRequest struct {
	A int `json:"a"`
	B int `json:"b"`
}

// SumResponse collects the response parameters for the Sum method.
type SumResponse struct {
	I0 int   `json:"i0"`
	E1 error `json:"e1"`
}

// MakeSumEndpoint returns an endpoint that invokes Sum on the service.
func MakeSumEndpoint(s service.HelloService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SumRequest)
		i0, e1 := s.Sum(ctx, req.A, req.B)
		return SumResponse{
			E1: e1,
			I0: i0,
		}, nil
	}
}

// Failed implements Failer.
func (r SumResponse) Failed() error {
	return r.E1
}

// ConcatRequest collects the request parameters for the Concat method.
type ConcatRequest struct {
	A string `json:"a"`
	B string `json:"b"`
}

// ConcatResponse collects the response parameters for the Concat method.
type ConcatResponse struct {
	S0 string `json:"s0"`
	E1 error  `json:"e1"`
}

// MakeConcatEndpoint returns an endpoint that invokes Concat on the service.
func MakeConcatEndpoint(s service.HelloService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ConcatRequest)
		s0, e1 := s.Concat(ctx, req.A, req.B)
		return ConcatResponse{
			E1: e1,
			S0: s0,
		}, nil
	}
}

// Failed implements Failer.
func (r ConcatResponse) Failed() error {
	return r.E1
}

// Sum implements Service. Primarily useful in a client.
func (e Endpoints) Sum(ctx context.Context, a int, b int) (i0 int, e1 error) {
	request := SumRequest{
		A: a,
		B: b,
	}
	response, err := e.SumEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(SumResponse).I0, response.(SumResponse).E1
}

// Concat implements Service. Primarily useful in a client.
func (e Endpoints) Concat(ctx context.Context, a string, b string) (s0 string, e1 error) {
	request := ConcatRequest{
		A: a,
		B: b,
	}
	response, err := e.ConcatEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(ConcatResponse).S0, response.(ConcatResponse).E1
}
