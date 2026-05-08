package config

import (
	"fmt"
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
func Load() (*Config, error) {
	// Загружаем .env файл (если есть), игнорируем ошибку для продакшена.
	if err := godotenv.Overload(); err != nil {
		return nil, fmt.Errorf("error read env file: %w", err)
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

	// Простейшая проверка обязательных SMTP-настроек
	if cfg.SMTP.Host == "" ||
		cfg.SMTP.Port == "" ||
		cfg.SMTP.Username == "" ||
		cfg.SMTP.Password == "" ||
		cfg.SMTP.From == "" ||
		cfg.SMTP.To == "" {
		return nil, fmt.Errorf("check config for SMTP vars in env")
	}

	return cfg, nil
}

// getEnv возвращает значение переменной окружения или значение по умолчанию.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
