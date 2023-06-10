package config

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"

	"formal.com/go-reverse-proxy-homework/internal/logging"
)

type Configuration struct {
	SourcePath         string   `json:"sourcePath"`
	TargetURL          string   `json:"targetURL"`
	BlockedHeaders     []string `json:"blockedHeaders"`
	BlockedQueryParams []string `json:"blockedQueryParams"`
}

type Configurations struct {
	Configs []Configuration `json:"configs"`
}

/**
 * PrintConfigurations prints the configurations.
 */
func (c *Configurations) PrintConfigurations(ctx context.Context) {
	logger := logging.GetLogger(ctx)

	for _, cfg := range c.Configs {
		logger.Println("Source Path:", cfg.SourcePath)
		logger.Println("Target URL:", cfg.TargetURL)
		logger.Println("Blocked Headers:", cfg.BlockedHeaders)
		logger.Println("Blocked Query Params:", cfg.BlockedQueryParams)
		logger.Println()
	}
}

func ParseConfig(r io.Reader) Configurations {
	config := Configurations{}
	err := json.NewDecoder(r).Decode(&config)
	if err != nil {
		log.Fatal("Failed to unmarshal config file: ", err)
	}

	return config
}

/**
 * ParseConfigFile parses the config file and converts it to the Configurations struct.
 * Requirement 5: You should provide a configuration file that allows the startup to define the rules for
 *		blocking requests.
 */
func ParseConfigFile(path string) Configurations {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal("Failed to open config file: ", err)
	}
	defer f.Close()

	return ParseConfig(f)
}
