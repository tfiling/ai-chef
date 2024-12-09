package configs

import (
	"os"
	"sync"
)

type Config struct {
	ClaudeAPIKey string
	MongoDBURI   string
}

var (
	instance Config
	once     sync.Once
)

func loadConfig() {
	instance = Config{
		ClaudeAPIKey: os.Getenv("CLAUDE_API_KEY"),
		MongoDBURI:   os.Getenv("MONGODB_URI"),
	}
}
func GetConfig() Config {
	once.Do(loadConfig)
	return instance
}
