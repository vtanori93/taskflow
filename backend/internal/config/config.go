package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config contiene la configuración de la aplicación
type Config struct {
	// Database
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string

	// Server
	ServerPort int
	ServerEnv  string

	// JWT
	JWTSecret            string
	JWTExpirationTime    int64
	JWTRefreshExpiration int64
}

// Load carga la configuración desde variables de entorno
func Load() (*Config, error) {
	// Cargar archivo .env si existe
	godotenv.Load()

	cfg := &Config{
		DBHost:               getEnv("POSTGRES_HOST", "localhost"),
		DBPort:               getEnvInt("POSTGRES_PORT", 5432),
		DBUser:               getEnv("POSTGRES_USER", "postgres"),
		DBPassword:           getEnv("POSTGRES_PASSWORD", "postgres"),
		DBName:               getEnv("POSTGRES_DB", "taskflow"),
		ServerPort:           getEnvInt("SERVER_PORT", 8080),
		ServerEnv:            getEnv("ENV", "development"),
		JWTSecret:            getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		JWTExpirationTime:    getEnvInt64("JWT_EXPIRATION_TIME", 3600),
		JWTRefreshExpiration: getEnvInt64("JWT_REFRESH_EXPIRATION", 604800),
	}

	return cfg, nil
}

// GetDSN retorna el DSN para conectar a PostgreSQL
func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
	)
}

// getEnv obtiene una variable de entorno con valor por defecto
func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}

// getEnvInt obtiene una variable de entorno como int
func getEnvInt(key string, defaultVal int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultVal
}

// getEnvInt64 obtiene una variable de entorno como int64
func getEnvInt64(key string, defaultVal int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intVal
		}
	}
	return defaultVal
}
