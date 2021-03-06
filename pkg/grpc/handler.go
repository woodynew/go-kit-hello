package grpc

import (
	"context"
	"errors"

	endpoint "github.com/woodynew/go-kit-hello/pkg/endpoint"
	pb "github.com/woodynew/go-kit-hello/pkg/grpc/pb"

	grpc "github.com/go-kit/kit/transport/grpc"
)

// makeFooHandler creates the handler logic
func makeFooHandler(endpoints endpoint.Endpoints, options []grpc.ServerOption) grpc.Handler {
	return grpc.NewServer(endpoints.FooEndpoint, decodeFooRequest, encodeFooResponse, options...)
}

// decodeFooResponse is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain Foo request.
// TODO implement the decoder
func decodeFooRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.FooRequest)
	return endpoint.FooRequest{S: req.S}, nil
	// return nil, errors.New("'Hello' Decoder is not impelemented")
}

// encodeFooResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
// TODO implement the encoder
func encodeFooResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(endpoint.FooResponse)
	return &pb.FooReply{S: string(resp.Rs), Err: err2str(resp.Err)}, nil
	// return nil, errors.New("'Hello' Encoder is not impelemented")
}
func (g *grpcServer) Foo(ctx context.Context, req *pb.FooRequest) (*pb.FooReply, error) {
	_, rep, err := g.foo.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.FooReply), nil
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
