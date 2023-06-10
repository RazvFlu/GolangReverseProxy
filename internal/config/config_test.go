package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseConfig(t *testing.T) {
	reader := strings.NewReader("{\"configs\": [  {\"sourcePath\": " +
		"\"test1\",\"targetURL\": \"http://test:8081\",\"blockedHeaders\": " +
		"[\"Authorization\"],\"blockedQueryParams\": " +
		"[\"password\"]  }]}")

	config := ParseConfig(reader)

	assert.Equal(t, 1, len(config.Configs))
	assert.Equal(t, "test1", config.Configs[0].SourcePath)
	assert.Equal(t, "http://test:8081", config.Configs[0].TargetURL)
	assert.Equal(t, 1, len(config.Configs[0].BlockedHeaders))
	assert.Equal(t, "Authorization", config.Configs[0].BlockedHeaders[0])
	assert.Equal(t, 1, len(config.Configs[0].BlockedQueryParams))
	assert.Equal(t, "password", config.Configs[0].BlockedQueryParams[0])
}
