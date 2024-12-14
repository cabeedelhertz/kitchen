package cors

import (
	"net/http"
	"time"

	"github.com/rs/cors"
)

// New create a new CORS configuration
func New(opts ...Option) *cors.Cors {
	return cors.New(parseOptions(opts...))
}

// parseOptions parses the set of cors options, defaulting as necessary
func parseOptions(opts ...Option) cors.Options {
	options := cors.Options{
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowOriginFunc: func(origin string) bool {
			// Allow all origins, which effectively disables CORS.
			return true
		},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{
			// Content-Type is in the default safelist.
			"Accept",
			"Accept-Encoding",
			"Accept-Post",
			"Connect-Accept-Encoding",
			"Connect-Content-Encoding",
			"Content-Encoding",
			"Grpc-Accept-Encoding",
			"Grpc-Encoding",
			"Grpc-Message",
			"Grpc-Status",
			"Grpc-Status-Details-Bin",
		},
		// Let browsers cache CORS information for longer, which reduces the
		// number of preflight requests. Any changes to ExposedHeaders won't
		// take effect until the cached data expires. FF caps this value at 24h,
		// and modern Chrome caps it at 2h.
		MaxAge: int(2 * time.Hour / time.Second),
	}
	for _, opt := range opts {
		opt(&options)
	}
	return options
}
