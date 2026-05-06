package mailer

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
	To       string // получатель заявок
}

func LoadConfig() Config {
	return Config{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
		Username: os.Getenv("SMTP_USER"),
		Password: os.Getenv("SMTP_PASSWORD"),
		From:     os.Getenv("SMTP_FROM"),
		To:       os.Getenv("SMTP_TO"),
	}
}

func (c Config) Send(subject, body string) error {
	if c.Host == "" || c.Port == "" || c.Username == "" || c.Password == "" || c.To == "" {
		return fmt.Errorf("SMTP configuration incomplete")
	}

	addr := fmt.Sprintf("%s:%s", c.Host, c.Port)

	// Формируем простое письмо в формате plain text
	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		c.From, c.To, subject, body,
	)

	auth := smtp.PlainAuth("", c.Username, c.Password, c.Host)

	// Отправляем через STARTTLS (порт 587)
	return smtp.SendMail(addr, auth, c.From, []string{c.To}, []byte(msg))
}

// Формирование текста письма из данных формы
func BuildContactBody(name, company, email, phone, message string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Имя: %s\n", name))
	if company != "" {
		sb.WriteString(fmt.Sprintf("Компания: %s\n", company))
	}
	sb.WriteString(fmt.Sprintf("Email: %s\n", email))
	if phone != "" {
		sb.WriteString(fmt.Sprintf("Телефон: %s\n", phone))
	}
	sb.WriteString(fmt.Sprintf("\nСообщение:\n%s\n", message))
	return sb.String()
}
