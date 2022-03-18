package thrift

import (
	"context"
	"errors"
	endpoint "hello/pkg/endpoint"
	addthrift "hello/pkg/thrift/gen-go/addsvc"
)

func (s *thriftServer) Foo(ctx context.Context, a string) (*addthrift.FooReply, error) {
	request := endpoint.FooRequest{S: string(a)}
	response, err := s.endpoints.FooEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	resp := response.(endpoint.FooResponse)
	return &addthrift.FooReply{Value: string(resp.Rs), Err: err2str(resp.Err)}, nil
}

func (s *thriftServer) Sum(ctx context.Context, a int64, b int64) (*addthrift.SumReply, error) {
	request := endpoint.SumRequest{A: int(a), B: int(b)}
	response, err := s.endpoints.SumEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	resp := response.(endpoint.SumResponse)
	return &addthrift.SumReply{Value: int64(resp.I0), Err: err2str(resp.E1)}, nil
}

func (s *thriftServer) Concat(ctx context.Context, a string, b string) (*addthrift.ConcatReply, error) {
	request := endpoint.ConcatRequest{A: a, B: b}
	response, err := s.endpoints.ConcatEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	resp := response.(endpoint.ConcatResponse)
	return &addthrift.ConcatReply{Value: resp.S0, Err: err2str(resp.E1)}, nil
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
