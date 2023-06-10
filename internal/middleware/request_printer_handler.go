package middleware

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"formal.com/go-reverse-proxy-homework/internal/logging"
	"formal.com/go-reverse-proxy-homework/internal/middleware/dynamic"
)

type printResponseWriter struct {
	http.ResponseWriter
	recorder *httptest.ResponseRecorder
}

func (prw *printResponseWriter) Write(data []byte) (int, error) {
	prw.recorder.Write(data)
	return prw.ResponseWriter.Write(data)
}

func (prw *printResponseWriter) WriteHeader(status int) {
	prw.recorder.WriteHeader(status)
	prw.ResponseWriter.WriteHeader(status)
}

func printRequest(requestPrefix string, responsePrefix string) dynamic.Handler {
	return func(ctx context.Context, rw http.ResponseWriter, req *http.Request, next dynamic.Next) {
		logger := logging.GetLogger(ctx)

		requestBody, _ := ioutil.ReadAll(req.Body)

		logger.Println(requestPrefix, "Request Headers:", req.Header)
		logger.Println(requestPrefix, "Request Body:", string(requestBody))

		req.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))

		recorder := httptest.NewRecorder()
		next(ctx, &printResponseWriter{ResponseWriter: rw, recorder: recorder}, req)

		logger.Println(responsePrefix, "Response Headers:", recorder.Header())
		logger.Println(responsePrefix, "Response Body:", recorder.Body)
	}
}

/**
 * PrintRequestHandler is a middleware that prints the request and response bodies and headers.
 * Requirement 3: The proxy should log all incoming requests, including headers and body, and response
 * 		headers and body.
 */
func PrintRequestHandler(requestPrefix string, responsePrefix string) dynamic.Handler {
	return func(ctx context.Context, rw http.ResponseWriter, req *http.Request, next dynamic.Next) {
		dynamic.Wrap(next, printRequest(requestPrefix, responsePrefix))(ctx, rw, req)
	}
}
