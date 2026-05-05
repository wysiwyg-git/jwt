package main

import (
	"log"
	"net/http"
	"os"

	"company-site/models"
	"company-site/templates"

	"github.com/a-h/templ"
)

func main() {
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
	log.Fatal(http.ListenAndServe("127.0.0.1:"+port, nil))
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
		success := r.URL.Query().Get("success") == "1"
		data := models.ContactFormData{}
		render(w, r, templates.ContactPage(data, success))
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
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
			// Успешно: здесь будет отправка письма или сохранение в БД
			log.Printf("Получено сообщение от %s (%s): %s", data.Name, data.Email, data.Message)
			// Перенаправление на страницу "Спасибо" или отображение success-сообщения
			// Воспользуемся query-параметром для простоты
			http.Redirect(w, r, "/contact?success=1", http.StatusSeeOther)
			return
		}
		// Есть ошибки — снова рендерим форму с данными и ошибками
		success := r.URL.Query().Get("success") == "1"
		render(w, r, templates.ContactPage(data, success))
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
