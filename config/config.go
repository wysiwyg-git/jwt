package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config хранит все конфигурационные параметры приложения.
type Config struct {
	ServerPort string
	SMTP       SMTPConfig
}

// SMTPConfig хранит настройки для отправки почты.
type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
	To       string
}

// Load загружает конфигурацию из .env файла и переменных окружения.
func Load() *Config {
	// Загружаем .env файл (если есть), игнорируем ошибку для продакшена.
	if err := godotenv.Overload(); err != nil {
		log.Println("Info: .env file not found, using system environment variables")
	}

	cfg := &Config{
		ServerPort: getEnv("PORT", "8080"),
		SMTP: SMTPConfig{
			Host:     os.Getenv("SMTP_HOST"),
			Port:     os.Getenv("SMTP_PORT"),
			Username: os.Getenv("SMTP_USER"),
			Password: os.Getenv("SMTP_PASSWORD"),
			From:     os.Getenv("SMTP_FROM"),
			To:       os.Getenv("SMTP_TO"),
		},
	}

	// Простейшая проверка обязательных SMTP-настроек (опционально)
	if cfg.SMTP.Host == "" || cfg.SMTP.Port == "" || cfg.SMTP.To == "" {
		log.Println("Warning: SMTP is not fully configured. Email sending will fail.")
	}

	return cfg
}

// getEnv возвращает значение переменной окружения или значение по умолчанию.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
