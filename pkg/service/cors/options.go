package cors

import "github.com/rs/cors"

// Option is a cors option
type Option func(*cors.Options)

// WithAllowedMethods sets the set of allowed CORS methods
func WithAllowedMethods(allowedMethods ...string) Option {
	return func(options *cors.Options) {
		options.AllowedMethods = allowedMethods
	}
}

// WithAllowedOriginsFunc sets the function to use to determine if an origin is
// allowed
func WithAllowedOriginsFunc(allowedOriginsFunc func(string) bool) Option {
	return func(options *cors.Options) {
		options.AllowOriginFunc = allowedOriginsFunc
	}
}

// WithAllowedHeaders sets the set of allowed CORS headers
func WithAllowedHeaders(allowedHeaders ...string) Option {
	return func(options *cors.Options) {
		options.AllowedHeaders = allowedHeaders
	}
}

// WithExposedHeaders sets the set of exposed headers
func WithExposedHeaders(exposedHeaders ...string) Option {
	return func(options *cors.Options) {
		options.ExposedHeaders = exposedHeaders
	}
}

// WithOptionsPassthrough instructs preflight to let other potential next
// handlers to process the OPTIONS method. Turn this on if your application
// handles OPTIONS.
func WithOptionsPassthrough(optionsPassthrough bool) Option {
	return func(options *cors.Options) {
		options.OptionsPassthrough = optionsPassthrough
	}
}

// WithMaxAge sets the CORS max age
func WithMaxAge(maxAge int) Option {
	return func(options *cors.Options) {
		options.MaxAge = maxAge
	}
}

// WithAllowCredentials indicates whether the request can include user
// credentials like cookies, HTTP authentication or client side SSL certificates
func WithAllowCredentials(allowCredentials bool) Option {
	return func(options *cors.Options) {
		options.AllowCredentials = allowCredentials
	}
}
