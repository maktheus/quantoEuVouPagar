package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config armazena as configurações da aplicação.
type Config struct {
	ServerPort string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
}

// LoadConfig carrega as configurações a partir de variáveis de ambiente.
func LoadConfig() *Config {
	// Determina qual arquivo .env carregar com base no ambiente
	env := os.Getenv("ENV")
	if env == "production" {
		godotenv.Load(".env.production")
	} else {
		// Em desenvolvimento, tenta carregar .env primeiro, depois .env.local
		godotenv.Load(".env", ".env.local")
	}

	return &Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "admin"),
		DBPassword: getEnv("DB_PASSWORD", "adminpassword"),
		DBName:     getEnv("DB_NAME", "quanto_eu_vou_pagar"),
		JWTSecret:  getEnv("JWT_SECRET_KEY", "my_secret_key"),
	}
}

// getEnv retorna o valor de uma variável de ambiente ou um valor padrão se ela não estiver definida.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}