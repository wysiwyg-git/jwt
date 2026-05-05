package models

import "strings"

type ContactFormData struct {
	Name    string
	Company string
	Email   string
	Phone   string
	Message string
	Errors  map[string]string // ключ — имя поля, значение — текст ошибки
}

func ValidateContactForm(data *ContactFormData) bool {
	if data.Errors == nil {
		data.Errors = make(map[string]string)
	}
	if strings.TrimSpace(data.Name) == "" {
		data.Errors["name"] = "Имя обязательно"
	}
	if strings.TrimSpace(data.Email) == "" {
		data.Errors["email"] = "Email обязателен"
	} else if !strings.Contains(data.Email, "@") || !strings.Contains(strings.Split(data.Email, "@")[1], ".") {
		data.Errors["email"] = "Некорректный email"
	}
	if strings.TrimSpace(data.Message) == "" {
		data.Errors["message"] = "Сообщение обязательно"
	}
	return len(data.Errors) == 0
}
