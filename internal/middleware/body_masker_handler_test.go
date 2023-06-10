package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"formal.com/go-reverse-proxy-homework/internal/middleware/dynamic"
	"github.com/stretchr/testify/assert"
)

func TestBodyMaskerHandlerWithoutSensitiveData(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	handler := dynamic.AsHttp(MaskResponseBodyHandler())

	handler(res, req, http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("Hello World!"))
	}))

	assert.Equal(t, "Hello World!", res.Body.String())
}

func TestBodyMaskerHandlerWithSensitiveData(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	handler := dynamic.AsHttp(MaskResponseBodyHandler())

	handler(res, req, http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("This is an example email: john.doe@example.com"))
	}))

	assert.Equal(t, "This is an example email: ********@example.com", res.Body.String())
}
