// Package config handles application configuration loading and validation.
// It reads settings from environment variables with sensible defaults.
package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all application configuration values.
type Config struct {
	// Server settings
	Host string
	Port int

	// Synology DiskStation settings
	DSHost     string
	DSPort     int
	DSUser     string
	DSPassword string
	DSHTTPS    bool

	// API settings
	APIKey     string
	DebugMode  bool
}

// Load reads configuration from environment variables and returns a Config.
// Returns an error if required fields are missing or invalid.
func Load() (*Config, error) {
	cfg := &Config{
		Host:      getEnv("HOST", "0.0.0.0"),
		Port:      getEnvInt("PORT", 8080),
		DSHost:    getEnv("DS_HOST", ""),
		DSPort:    getEnvInt("DS_PORT", 5000),
		DSUser:    getEnv("DS_USER", ""),
		DSPassword: getEnv("DS_PASSWORD", ""),
		DSHTTPS:   getEnvBool("DS_HTTPS", false),
		APIKey:    getEnv("API_KEY", ""),
		DebugMode: getEnvBool("DEBUG", false),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validate checks that all required configuration fields are set.
func (c *Config) validate() error {
	if c.DSHost == "" {
		return fmt.Errorf("DS_HOST is required")
	}
	if c.DSUser == "" {
		return fmt.Errorf("DS_USER is required")
	}
	if c.DSPassword == "" {
		return fmt.Errorf("DS_PASSWORD is required")
	}
	return nil
}

// DSBaseURL returns the base URL for the Synology DiskStation API.
func (c *Config) DSBaseURL() string {
	scheme := "http"
	if c.DSHTTPS {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s:%d", scheme, c.DSHost, c.DSPort)
}

// ServerAddr returns the formatted listen address for the HTTP server.
func (c *Config) ServerAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// getEnv retrieves a string environment variable or returns a default value.
func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

// getEnvInt retrieves an integer environment variable or returns a default value.
func getEnvInt(key string, defaultVal int) int {
	if val, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}

// getEnvBool retrieves a boolean environment variable or returns a default value.
func getEnvBool(key string, defaultVal bool) bool {
	if val, ok := os.LookupEnv(key); ok {
		if b, err := strconv.ParseBool(val); err == nil {
			return b
		}
	}
	return defaultVal
}
