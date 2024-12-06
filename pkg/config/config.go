package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Host            string
	Port            string
	Dsn             string
	MaxSendFileSize int64
	DialerTimeout   time.Duration
}

// New returns a new Config struct
func New() *Config {
	return &Config{
		Host:            getEnv("HOST", "localhost"),
		Port:            getEnv("PORT", "4000"),
		Dsn:             getEnv("DSN", "./forum.db"),
		MaxSendFileSize: int64(getEnvAsInt("MAX_SEND_FILE_SIZE", 26214400)),
		DialerTimeout:   time.Duration(getEnvAsInt("DIALER_TIMEOUT", 60)),
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

// Simple helper function to read an environment variable into integer or return a default value
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

// Helper to read an environment variable into a bool or return default value
func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}

// Helper to read an environment variable into a string slice or return default value
func getEnvAsSlice(name string, defaultVal []string, sep string) []string {
	valStr := getEnv(name, "")

	if valStr == "" {
		return defaultVal
	}

	val := strings.Split(valStr, sep)

	return val
}
