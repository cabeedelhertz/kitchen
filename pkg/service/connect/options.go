package connect

import (
	"kitchen/pkg/service/cors"

	"connectrpc.com/connect"
)

// Option is a configuration option
type Option func(*options)

// options configuration options
type options struct {
	additionalServices []ServiceRegistrar
	handlerOptions     []connect.HandlerOption
	corsOptions        []cors.Option
}

// WithAdditionalService specifies an additional service to mount on the connect
// server
func WithAdditionalService[T any](factory HandlerFactory[T], handler T) Option {
	return func(options *options) {
		options.additionalServices = append(options.additionalServices, Service[T]{factory, handler})
	}
}

// WithHandlerOptions passes the set of connect options to use
func WithHandlerOptions(handlerOptions ...connect.HandlerOption) Option {
	return func(options *options) {
		options.handlerOptions = handlerOptions
	}
}

// WithCORSOptions passes the set of cors options to use
func WithCORSOptions(corsOptions ...cors.Option) Option {
	return func(options *options) {
		options.corsOptions = corsOptions
	}
}
