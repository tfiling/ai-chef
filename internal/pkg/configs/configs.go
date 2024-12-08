package configs

import (
	"os"
	"sync"
)

type Config struct {
	ClaudeAPIKey string `mapstructure:"api_key"`
}

var (
	instance Config
	once     sync.Once
)

func loadConfig() {
	instance = Config{ClaudeAPIKey: os.Getenv("CLAUDE_API_KEY")}
}
func GetConfig() Config {
	once.Do(loadConfig)
	return instance
}
