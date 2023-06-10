package middleware

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"formal.com/go-reverse-proxy-homework/internal/config"
	"formal.com/go-reverse-proxy-homework/internal/logging"
	"formal.com/go-reverse-proxy-homework/internal/middleware/dynamic"
	"github.com/stretchr/testify/assert"
)

func TestRequestBlockerWhenRequestIsNotBlocked(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	config := config.Configuration{
		SourcePath:         "/test",
		TargetURL:          "http://test:8081",
		BlockedHeaders:     []string{},
		BlockedQueryParams: []string{},
	}

	handler := dynamic.AsHttp(BlockRequestHandler(config))

	handler(res, req, http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("Hello World!"))
	}))

	assert.Equal(t, "Hello World!", res.Body.String())
	assert.Equal(t, http.StatusOK, res.Code)
}

func TestRequestBlockerWhenRequestIsBlockedByHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	logger := log.New(os.Stdout, "", 0)

	ctx := context.WithValue(req.Context(), logging.LoggerKey, logger)
	req = req.WithContext(ctx)

	config := config.Configuration{
		SourcePath:         "/test",
		TargetURL:          "http://test:8081",
		BlockedHeaders:     []string{"BlockedHeader"},
		BlockedQueryParams: []string{},
	}

	req.Header.Add("BlockedHeader", "BlockedHeader")

	handler := dynamic.AsHttp(BlockRequestHandler(config))

	handler(res, req, http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("Hello World!"))
	}))

	assert.Equal(t, http.StatusForbidden, res.Code)
}

func TestRequestBlockerWhenRequestIsBlockedByQueryParam(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test?BlockedParam=BlockedParam", nil)
	res := httptest.NewRecorder()
	logger := log.New(os.Stdout, "", 0)

	ctx := context.WithValue(req.Context(), logging.LoggerKey, logger)
	req = req.WithContext(ctx)

	config := config.Configuration{
		SourcePath:         "/test",
		TargetURL:          "http://test:8081",
		BlockedHeaders:     []string{},
		BlockedQueryParams: []string{"BlockedParam"},
	}

	handler := dynamic.AsHttp(BlockRequestHandler(config))

	handler(res, req, http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("Hello World!"))
	}))

	assert.Equal(t, http.StatusForbidden, res.Code)
}
