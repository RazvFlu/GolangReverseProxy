package logging

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLogger(t *testing.T) {
	ctx := context.Background()
	logger := log.New(nil, "", 0)

	ctx = context.WithValue(ctx, LoggerKey, logger)
	contextLogger := GetLogger(ctx)

	assert.Equal(t, logger, contextLogger)
}
