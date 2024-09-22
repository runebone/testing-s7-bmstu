package config

import (
	"os"
	"strconv"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Database DatabaseConfig
	Logger   LoggerConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type LoggerConfig struct {
	Level string
}

func LoadConfig(path string) (Config, error) {
	var config Config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return config, err
	}

	config.Database.Host = getEnv("DATABASE_HOST", config.Database.Host)
	config.Database.Port = getEnvInt("DATABASE_PORT", config.Database.Port)
	config.Database.User = getEnv("DATABASE_USER", config.Database.User)
	config.Database.Password = getEnv("DATABASE_PASSWORD", config.Database.Password)
	config.Database.DBName = getEnv("DATABASE_NAME", config.Database.DBName)
	config.Database.SSLMode = getEnv("DATABASE_SSLMODE", config.Database.SSLMode)

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
