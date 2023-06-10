package dynamic

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/negroni"
)

func TestDynamicServerMiddlewareCalls(t *testing.T) {
	numberOfCalls := 0

	testFunction := func() Handler {
		return func(ctx context.Context, rw http.ResponseWriter, req *http.Request, next Next) {
			numberOfCalls++
			next(ctx, rw, req)
		}
	}

	testHandler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request, next Next) {
		Wrap(next, testFunction())(ctx, rw, req)
	}

	serverChain := DynamicServerMiddleware(func(req *http.Request) (chain []negroni.Handler) {
		chain = append(chain, AsHttp(testHandler))
		chain = append(chain, AsHttp(testHandler))
		chain = append(chain, AsHttp(testHandler))

		return
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	serverChain(res, req, http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {}))
	assert.Equal(t, 3, numberOfCalls)
}
