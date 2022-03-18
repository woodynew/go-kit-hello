package grpc

import (
	"context"
	"errors"
	endpoint1 "hello/pkg/endpoint"
	pb "hello/pkg/grpc/pb"
	service "hello/pkg/service"

	endpoint "github.com/go-kit/kit/endpoint"
	grpc1 "github.com/go-kit/kit/transport/grpc"
	grpc "google.golang.org/grpc"
)

// New returns an AddService backed by a gRPC server at the other end
//  of the conn. The caller is responsible for constructing the conn, and
// eventually closing the underlying transport. We bake-in certain middlewares,
// implementing the client library pattern.
func New(conn *grpc.ClientConn, options map[string][]grpc1.ClientOption) (service.HelloService, error) {
	var fooEndpoint endpoint.Endpoint
	{
		fooEndpoint = grpc1.NewClient(conn, "pb.Hello", "Foo", encodeFooRequest, decodeFooResponse, pb.FooReply{}, options["Foo"]...).Endpoint()
	}

	return endpoint1.Endpoints{FooEndpoint: fooEndpoint}, nil
}

// encodeFooRequest is a transport/grpc.EncodeRequestFunc that converts a
//  user-domain Foo request to a gRPC request.
func encodeFooRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(endpoint1.FooRequest)
	return &pb.FooRequest{S: string(req.S)}, nil
	// return nil, errors.New("'Hello' Encoder is not impelemented")
}

// decodeFooResponse is a transport/grpc.DecodeResponseFunc that converts
// a gRPC concat reply to a user-domain concat response.
func decodeFooResponse(_ context.Context, reply interface{}) (interface{}, error) {
	rep := reply.(*pb.FooReply)
	return endpoint1.FooResponse{Rs: rep.S, Err: str2err(rep.Err)}, nil
	// return nil, errors.New("'Hello' Decoder is not impelemented-client")
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
