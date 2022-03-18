package service

import (
	"context"
	"fmt"

	log "github.com/go-kit/kit/log"
)

// Middleware describes a service middleware.
type Middleware func(HelloService) HelloService

type hiMiddleware struct {
	next HelloService
}

func HiMiddleware() Middleware {
	return func(next HelloService) HelloService {
		return &hiMiddleware{next}
	}

}
func (h hiMiddleware) Foo(ctx context.Context, s string) (rs string, err error) {

	return h.next.Foo(ctx, s)
}
func (h hiMiddleware) Sum(ctx context.Context, a int, b int) (i0 int, e1 error) {

	return h.next.Sum(ctx, a, b)
}
func (h hiMiddleware) Concat(ctx context.Context, a string, b string) (s0 string, e1 error) {

	return h.next.Concat(ctx, a, b)
}

type loggingMiddleware struct {
	logger log.Logger
	next   HelloService
}

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next HelloService) HelloService {
		return &loggingMiddleware{logger, next}
	}

}

func (l loggingMiddleware) Foo(ctx context.Context, s string) (rs string, err error) {
	defer func() {
		l.logger.Log("method", "Foo", "s", s, "rs", rs, "err", err)
	}()
	return l.next.Foo(ctx, s)
}
func (l loggingMiddleware) Sum(ctx context.Context, a int, b int) (i0 int, e1 error) {
	defer func() {
		l.logger.Log("method", "Sum", "a", a, "b", b, "i0", i0, "e1", e1)
	}()
	fmt.Println("------loggingMiddleware-Sum----------")
	return l.next.Sum(ctx, a, b)
}
func (l loggingMiddleware) Concat(ctx context.Context, a string, b string) (s0 string, e1 error) {
	defer func() {
		l.logger.Log("method", "Concat", "a", a, "b", b, "s0", s0, "e1", e1)
	}()
	return l.next.Concat(ctx, a, b)
}
