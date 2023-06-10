package proxy

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"formal.com/go-reverse-proxy-homework/internal/config"
	"formal.com/go-reverse-proxy-homework/internal/logging"
	"github.com/stretchr/testify/assert"
)

func TestProxyWithoutProxyBackend(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", bytes.NewBuffer([]byte("Hello World!")))
	res := httptest.NewRecorder()
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)

	ctx := context.WithValue(req.Context(), logging.LoggerKey, logger)
	req = req.WithContext(ctx)

	config := config.Configuration{
		SourcePath:         "/test",
		TargetURL:          "http://test:8081",
		BlockedHeaders:     []string{},
		BlockedQueryParams: []string{},
	}

	proxy := NewProxy(ctx, config)

	proxy.ServeHTTP(res, req)

	assert.Equal(t, http.StatusInternalServerError, res.Code)
}

func TestInspectableRequest(t *testing.T) {
	getReq := httptest.NewRequest(http.MethodGet, "/", nil)
	postReq := httptest.NewRequest(http.MethodPost, "/", nil)
	putReq := httptest.NewRequest(http.MethodPut, "/", nil)
	patchReq := httptest.NewRequest(http.MethodPatch, "/", nil)
	deleteReq := httptest.NewRequest(http.MethodDelete, "/", nil)

	assert.True(t, isInspectableRequest(getReq))
	assert.False(t, isInspectableRequest(postReq))
	assert.False(t, isInspectableRequest(putReq))
	assert.False(t, isInspectableRequest(patchReq))
	assert.False(t, isInspectableRequest(deleteReq))
}
