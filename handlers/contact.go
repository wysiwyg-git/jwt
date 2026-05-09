package handlers

import (
	"company-site/config"
	"company-site/mailer"
	"company-site/models"
	"company-site/templates"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

var (
	submissionsMu sync.Mutex
	submissions   = make(map[string]time.Time)
)

// ContactHandler – страница с контактной формой
func ContactHandler(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	switch r.Method {
	case http.MethodGet:
		success := r.URL.Query().Get("success")
		var data models.ContactFormData
		render(w, r, templates.ContactPage(data, success))
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// rate limiting: не чаще раза в минуту с одного IP
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr // fallback
		}

		submissionsMu.Lock()
		last, exists := submissions[ip]
		submissionsMu.Unlock()
		if exists && time.Since(last) < time.Minute {
			http.Error(w, "Слишком много запросов. Пожалуйста, подождите 1 минуту.", http.StatusTooManyRequests)
			return
		}

		data := models.ContactFormData{
			Name:    r.FormValue("name"),
			Company: r.FormValue("company"),
			Email:   r.FormValue("email"),
			Phone:   r.FormValue("phone"),
			Message: r.FormValue("message"),
		}
		if models.ValidateContactForm(&data) {
			// Отправка письма
			body := mailer.BuildContactBody(data.Name, data.Company, data.Email, data.Phone, data.Message)
			if err := mailer.Send(cfg.SMTP, "Новая заявка с сайта ПромКлей", body); err != nil {
				log.Printf("Ошибка отправки email: %v", err)
				http.Redirect(w, r, "/contact?success=0", http.StatusSeeOther)
				return
			}

			submissionsMu.Lock()
			submissions[ip] = time.Now()
			submissionsMu.Unlock()

			// Успех
			http.Redirect(w, r, "/contact?success=1", http.StatusSeeOther)
			return
		}
		// Ошибки валидации – показываем форму снова
		render(w, r, templates.ContactPage(data, ""))
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
