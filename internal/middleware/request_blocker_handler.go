package middleware

import (
	"context"
	"net/http"

	"formal.com/go-reverse-proxy-homework/internal/config"
	"formal.com/go-reverse-proxy-homework/internal/logging"
	"formal.com/go-reverse-proxy-homework/internal/middleware/dynamic"
)

func isBlockedRequest(req *http.Request, config config.Configuration) bool {
	for _, header := range config.BlockedHeaders {
		if req.Header.Get(header) != "" {
			return true
		}
	}

	for _, queryParam := range config.BlockedQueryParams {
		if req.URL.Query().Get(queryParam) != "" {
			return true
		}
	}

	return false
}

func blockRequest(config config.Configuration) dynamic.Handler {
	return func(ctx context.Context, rw http.ResponseWriter, req *http.Request, next dynamic.Next) {
		logger := logging.GetLogger(ctx)

		if isBlockedRequest(req, config) {
			logger.Println("Blocked Request")
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		next(ctx, rw, req)
	}
}

/**
 * BlockRequestHandler is a middleware that blocks requests based on the configuration.
 * Requirement 4: The proxy should be able to block requests based on a set of predefined rules. For
 *		example, you can block requests that contain specific headers or parameters.
 */
func BlockRequestHandler(config config.Configuration) dynamic.Handler {
	return func(ctx context.Context, rw http.ResponseWriter, req *http.Request, next dynamic.Next) {
		dynamic.Wrap(next, blockRequest(config))(ctx, rw, req)
	}
}
