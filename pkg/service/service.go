package service

import "context"

// HelloService describes the service.
type HelloService interface {
	// Add your methods here
	Foo(ctx context.Context, s string) (rs string, err error)
	Sum(ctx context.Context, a, b int) (int, error)
	Concat(ctx context.Context, a, b string) (string, error)
}

type basicHelloService struct{}

func (b *basicHelloService) Foo(ctx context.Context, s string) (rs string, err error) {
	return "hhhhhhhaaaaaa", err
}

func (ba *basicHelloService) Sum(ctx context.Context, a int, b int) (i0 int, e1 error) {
	// TODO implement the business logic of Sum
	return a + b, e1
}
func (ba *basicHelloService) Concat(ctx context.Context, a string, b string) (s0 string, e1 error) {
	// TODO implement the business logic of Concat
	return a + b, e1
}

// NewBasicHelloService returns a naive, stateless implementation of HelloService.
func NewBasicHelloService() HelloService {
	return &basicHelloService{}
}

// New returns a HelloService with all of the expected middleware wired in.
func New(middleware []Middleware) HelloService {
	var svc HelloService = NewBasicHelloService()
	for _, m := range middleware {
		svc = m(svc)
	}
	return svc
}
