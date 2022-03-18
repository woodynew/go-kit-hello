package thrift

import (
	"context"
	endpoint "hello/pkg/endpoint"
	addthrift "hello/pkg/thrift/gen-go/addsvc"
)

type thriftServer struct {
	ctx       context.Context
	endpoints endpoint.Endpoints
}

// NewThriftServer makes a set of endpoints available as a Thrift service.
func NewThriftServer(endpoints endpoint.Endpoints) addthrift.AddService {
	return &thriftServer{
		endpoints: endpoints,
	}
}
