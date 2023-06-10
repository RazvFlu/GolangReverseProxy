package dynamic

import (
	"context"
	"net/http"

	"github.com/urfave/negroni"
)

type Next func(ctx context.Context, rw http.ResponseWriter, req *http.Request)
type Handler func(ctx context.Context, rw http.ResponseWriter, req *http.Request, next Next)

/**
 * DynamicServerMiddleware creates a new negroni middleware chain.
 *
 * For requests the chain is executed in the order they are added.
 * For responses the chain is executed in the reverse order they are added.
 */
func DynamicServerMiddleware(fn func(req *http.Request) []negroni.Handler) negroni.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		handlers := fn(req)
		chain := negroni.New(handlers...)
		chain.UseHandler(next)
		chain.ServeHTTP(rw, req)
	}
}

/**
 * Wrap creates a new dynamic.Next function which executes the given handlers in the order they are given.
 */
func Wrap(next Next, wares ...Handler) Next {
	final := next

	for i := len(wares) - 1; i >= 0; i-- {
		ware := wares[i]
		next := final
		final = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
			ware(ctx, rw, req, next)
		}
	}

	return final
}

/**
 * AsHttp converts a dynamic.Handler to a negroni.HandlerFunc.
 */
func AsHttp(handler Handler) negroni.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		handler(req.Context(), rw, req, func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
			next(rw, req)
		})
	}
}
