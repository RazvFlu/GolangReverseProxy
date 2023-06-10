package middleware

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"formal.com/go-reverse-proxy-homework/internal/logging"
	"formal.com/go-reverse-proxy-homework/internal/middleware/dynamic"
	"github.com/stretchr/testify/assert"
)

func TestRequestPrinter(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)

	ctx := context.WithValue(req.Context(), logging.LoggerKey, logger)
	req = req.WithContext(ctx)

	req.Header.Add("TestRequestHeader", "TestRequestHeader")
	req.Body = ioutil.NopCloser(bytes.NewBufferString("TestRequestBody"))

	handler := dynamic.AsHttp(PrintRequestHandler("requestPrefix", "responsePrefix"))

	handler(res, req, http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("TestResponseBody"))
	}))

	printerOutput := buf.String()

	log.Println(printerOutput)

	assert.Contains(t, printerOutput, "requestPrefix Request Headers: map[Testrequestheader:[TestRequestHeader]]")
	assert.Contains(t, printerOutput, "requestPrefix Request Body: TestRequestBody")
	assert.Contains(t, printerOutput, "responsePrefix Response Headers: map[")
	assert.Contains(t, printerOutput, "responsePrefix Response Body: TestResponseBody")
}
