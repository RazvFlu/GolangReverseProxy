package middleware

import (
	"context"
	"net/http"
	"regexp"
	"strings"

	"formal.com/go-reverse-proxy-homework/internal/middleware/dynamic"
)

type sensitiveInformationWriter struct {
	http.ResponseWriter
	Body []byte
}

/**
 * In reality we would want to use a more sophisticated approach to mask sensitive information
 * 		because response bodies can have big sizes and buffering it is not a good idea.
 * Maybe a sliding window approach with the chunks received here would be better.
 *  	You wound scan the body chunks multiple times but you would not lose conextual data and
 * 		you would not risk on receiving a chunk which splits the sensitive data.
 */
func (si *sensitiveInformationWriter) Write(b []byte) (int, error) {
	si.Body = append(si.Body, b...)
	return len(b), nil
}

func maskEmails(body []byte) []byte {
	regex := regexp.MustCompile(`(?i)\b[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}\b`)
	modifiedBody := regex.ReplaceAllFunc(body, func(match []byte) []byte {
		email := string(match)
		masked := maskEmail(email)
		return []byte(masked)
	})
	return modifiedBody
}

func maskEmail(email string) string {
	atIndex := strings.Index(email, "@")
	if atIndex == -1 {
		return email
	}

	username := email[:atIndex]
	maskedUsername := strings.Repeat("*", len(username))
	domain := email[atIndex:]
	return maskedUsername + domain
}

func modifyResponseBody(body []byte) []byte {
	// Use regex to replace content in the response body
	modifiedBody := maskEmails(body)

	return modifiedBody
}

func maskResponseBody() dynamic.Handler {
	return func(ctx context.Context, rw http.ResponseWriter, req *http.Request, next dynamic.Next) {
		sensitiveInformationModifier := &sensitiveInformationWriter{
			ResponseWriter: rw,
			Body:           nil,
		}

		next(ctx, sensitiveInformationModifier, req)

		// Modify the response body
		modifiedBody := modifyResponseBody(sensitiveInformationModifier.Body)
		rw.Write(modifiedBody)
	}
}

/**
 * MaskResponseBodyHandler is a middleware that masks sensitive information in the response body.
 * Requirement 2: The proxy should inspect the response body for any sensitive information and mask it
 *		before forwarding the response.
 */
func MaskResponseBodyHandler() dynamic.Handler {
	return func(ctx context.Context, rw http.ResponseWriter, req *http.Request, next dynamic.Next) {
		dynamic.Wrap(next, maskResponseBody())(ctx, rw, req)
	}
}
