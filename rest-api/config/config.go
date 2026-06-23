package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	MikroTik MikroTikConfig
	Server   ServerConfig
}

type MikroTikConfig struct {
	Host     string
	Username string
	Password string
	Timeout  time.Duration
	Retries  int
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func Load() (*Config, error) {
	retries, err := strconv.Atoi(getEnv("MIKROTIK_RETRIES", "3"))
	if err != nil {
		return nil, fmt.Errorf("invalid value MIKROTIK_RETRIES: %w", err)
	}

	return &Config{
		MikroTik: MikroTikConfig{
			Host:     getEnv("MIKROTIK_HOST", "192.168.88.1"),
			Username: getEnv("MIKROTIK_USER", "admin"),
			Password: mustEnv("MIKROTIK_PASS"),
			Timeout:  10 * time.Second,
			Retries:  retries,
		},
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
		},
	}, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		fmt.Fprintf(os.Stderr, "FATAL: environment variable %s not installed\n", key)
		os.Exit(1)
	}
	return v
}
