package main

import (
	"log"
	"net/http"
	"os"

	"company-site/templates"
)

func main() {
	// Раздача статических файлов из папки static/
	// http.Dir указывает корневую папку для файлового сервера.
	// StripPrefix убирает префикс /static/ из URL, чтобы
	// /static/css/bootstrap.min.css превратилось в css/bootstrap.min.css
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Обработчик главной страницы
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Если путь не корневой, отдаём 404 (чтобы избежать дублирования главной)
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		// Рендерим компонент IndexPage в ResponseWriter
		component := templates.IndexPage()
		if err := component.Render(r.Context(), w); err != nil {
			http.Error(w, "Ошибка рендеринга", http.StatusInternalServerError)
			log.Printf("Ошибка рендеринга: %v", err)
		}
	})

	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Сервер запущен на http://localhost:%s", port)
	log.Fatal(http.ListenAndServe("127.0.0.1:"+port, nil))
}
