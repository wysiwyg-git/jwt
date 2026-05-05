package main

import (
	"log"
	"net/http"
	"os"

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
	http.HandleFunc("/contact", func(w http.ResponseWriter, r *http.Request) {
		render(w, r, templates.ContactPage())
	})

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
