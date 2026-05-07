package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"company-site/mailer"
	"company-site/models"
	"company-site/templates"

	"github.com/a-h/templ"
	"github.com/joho/godotenv"
)

var (
	submissionsMu sync.Mutex
	submissions   = make(map[string]time.Time)
)

func main() {
	if err := godotenv.Overload(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Статические файлы
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// --- Маршруты страниц ---
	// Главная
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Только точный путь "/"
		if r.URL.Path != "/" {
			// Если путь не "/", отдаём 404
			renderNotFound(w, r)
			return
		}
		render(w, r, templates.IndexPage())
	})

	// Страница «О компании»
	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		render(w, r, templates.AboutPage())
	})

	// Услуги
	http.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
		render(w, r, templates.ServicesPage())
	})

	// Контакты
	http.HandleFunc("/contact", contactHandler)

	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Сервер запущен на http://localhost:%s", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}

// render рендерит любой templ.Component и обрабатывает ошибки
func render(w http.ResponseWriter, r *http.Request, comp templ.Component) {
	if err := comp.Render(r.Context(), w); err != nil {
		log.Printf("Ошибка рендеринга %s: %v", r.URL.Path, err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
	}
}

// renderNotFound показывает кастомную страницу 404
func renderNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	if err := templates.NotFoundPage().Render(r.Context(), w); err != nil {
		log.Printf("Ошибка рендеринга 404: %v", err)
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
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
			cfg := mailer.LoadConfig()
			body := mailer.BuildContactBody(data.Name, data.Company, data.Email, data.Phone, data.Message)
			if err := cfg.Send("Новая заявка с сайта ПромКлей", body); err != nil {
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
