package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	DB       DatabaseConfig
	Server   ServerConfig
	Cache    CacheConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type ServerConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type CacheConfig struct {
	TTL  time.Duration
	Size int
}

func LoadFromEnv() (Config, error) {
	cfg := Config{}

	cfg.DB.Host = getEnv("DB_HOST", "localhost")
	cfg.DB.Port = getEnvAsInt("DB_PORT", 3306)
	cfg.DB.User = getRequiredEnv("DB_USER")
	cfg.DB.Password = getRequiredEnv("DB_PASSWORD")
	cfg.DB.Name = getEnv("DB_NAME", "player_activity")

	cfg.Server.Port = getEnvAsInt("SERVER_PORT", 8080)
	cfg.Server.ReadTimeout = time.Duration(getEnvAsInt("SERVER_READ_TIMEOUT", 5)) * time.Second
	cfg.Server.WriteTimeout = time.Duration(getEnvAsInt("SERVER_WRITE_TIMEOUT", 10)) * time.Second

	cfg.Cache.TTL = time.Duration(getEnvAsInt("CACHE_TTL", 60)) * time.Minute
	cfg.Cache.Size = getEnvAsInt("CACHE_SIZE", 1000)

	if err := cfg.validate(); err != nil {
		return Config{}, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

func (c Config) validate() error {
	if c.DB.User == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if c.DB.Password == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("SERVER_PORT must be between 1 and 65535")
	}
	if c.Cache.Size < 1 {
		return fmt.Errorf("CACHE_SIZE must be greater than 0")
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getRequiredEnv(key string) string {
	return os.Getenv(key)
}

func getEnvAsInt(key string, defaultValue int) int {
	strValue := os.Getenv(key)
	if strValue == "" {
		return defaultValue
	}
	if value, err := strconv.Atoi(strValue); err == nil {
		return value
	}
	return defaultValue
}

func (c DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", c.User, c.Password, c.Host, c.Port, c.Name)
}