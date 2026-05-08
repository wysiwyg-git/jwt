package mailer

import (
	"company-site/config"
	"fmt"
	"net/smtp"
	"strings"
)

func Send(cfg config.SMTPConfig, subject, body string) error {
	if cfg.Host == "" || cfg.Port == "" || cfg.Username == "" || cfg.Password == "" || cfg.To == "" {
		return fmt.Errorf("SMTP configuration incomplete")
	}

	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	// Формируем простое письмо в формате plain text
	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		cfg.From, cfg.To, subject, body,
	)

	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)

	// Отправляем через STARTTLS (порт 587)
	return smtp.SendMail(addr, auth, cfg.From, []string{cfg.To}, []byte(msg))
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
