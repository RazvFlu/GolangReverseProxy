package proxy

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"formal.com/go-reverse-proxy-homework/internal/config"
	"formal.com/go-reverse-proxy-homework/internal/logging"
	"formal.com/go-reverse-proxy-homework/internal/middleware"
	"formal.com/go-reverse-proxy-homework/internal/middleware/dynamic"
	"github.com/urfave/negroni"
)

/**
 * A request is inspectable if it is a GET request.
 * Requirement 1: The reverse proxy should be able to forward any kind of request but inspect only GET requests.
 */
func isInspectableRequest(req *http.Request) bool {
	return req.Method == http.MethodGet
}

/**
 * Create a new proxy handler.
 * The DynamicServerMiddleware function is used to create a new negroni middleware chain
 * 		which works as a chain of responsibility where each handler does its work and calls the next handler.
 */
func NewProxy(ctx context.Context, config config.Configuration) http.HandlerFunc {
	backendUrl, err := url.Parse(config.TargetURL)
	if err != nil {
		log.Fatal("Failed to parse target URL: ", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(backendUrl)
	proxy.ErrorLog = log.New(logging.GetLogger(ctx).Writer(), "", 0)
	proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, err error) {
		log.Print("Proxy error: ", err)

		writer.WriteHeader(http.StatusInternalServerError)
	}

	serverChain := dynamic.DynamicServerMiddleware(func(req *http.Request) (chain []negroni.Handler) {
		chain = append(chain, dynamic.AsHttp(middleware.PrintRequestHandler("Before Processing:", "After Processing:")))
		chain = append(chain, dynamic.AsHttp(middleware.BlockRequestHandler(config)))

		if isInspectableRequest(req) {
			chain = append(chain, dynamic.AsHttp(middleware.MaskResponseBodyHandler()))
			chain = append(chain, dynamic.AsHttp(middleware.PrintRequestHandler("After Processing:", "Before Processing:")))
		}

		return
	})

	return func(rw http.ResponseWriter, req *http.Request) {
		serverChain(rw, req, proxy.ServeHTTP)
	}
}
