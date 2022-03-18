// THIS FILE IS AUTO GENERATED BY GK-CLI DO NOT EDIT!!
package grpc

import (
	grpc "github.com/go-kit/kit/transport/grpc"
	endpoint "github.com/woodynew/go-kit-hello/pkg/endpoint"
	pb "github.com/woodynew/go-kit-hello/pkg/grpc/pb"
)

// NewGRPCServer makes a set of endpoints available as a gRPC AddServer
type grpcServer struct {
	foo grpc.Handler
	pb.UnimplementedHelloServer
}

func NewGRPCServer(endpoints endpoint.Endpoints, options map[string][]grpc.ServerOption) pb.HelloServer {
	return &grpcServer{foo: makeFooHandler(endpoints, options["Foo"])}
}
