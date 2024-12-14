package metadata

import (
	"context"
	"kitchen/pkg/service/metadata"

	"connectrpc.com/connect"
)

var _ connect.Interceptor = (*Interceptor)(nil)

// Interceptor is a metadata interceptor that adds the request metadata to the
// context
type Interceptor struct{}

// NewInterceptor creates a new connect metadata interceptor
func NewInterceptor() *Interceptor {
	return new(Interceptor)
}

// WrapUnary wraps a unary call adding the incoming metadata to the context
func (i *Interceptor) WrapUnary(fn connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
		return fn(metadata.NewIncomingContext(ctx, metadata.MD(request.Header())), request)
	}
}

// WrapStreamingClient wraps a streaming client call adding the outgoing
// metadata to the  context
func (i *Interceptor) WrapStreamingClient(fn connect.StreamingClientFunc) connect.StreamingClientFunc {
	return func(ctx context.Context, spec connect.Spec) connect.StreamingClientConn {
		return fn(metadata.NewOutgoingContext(ctx, metadata.MD{}), spec)
	}
}

// WrapStreamingHandler wraps a streaming handler call adding the incoming
// metadata to the context
func (i *Interceptor) WrapStreamingHandler(fn connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		return fn(metadata.NewIncomingContext(ctx, metadata.MD(conn.RequestHeader())), conn)
	}
}
