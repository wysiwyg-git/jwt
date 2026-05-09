package handlers

import (
	"company-site/config"
	"net/http"
)

func SetupRoutes(mux *http.ServeMux, cfg *config.Config) {
	// Оборачиваем обработчики, которые нуждаются в конфиге
	mux.HandleFunc("/contact", func(w http.ResponseWriter, r *http.Request) {
		ContactHandler(w, r, cfg)
	})

	// Остальные обработчики без доп. зависимостей
	mux.HandleFunc("/", IndexHandler)
	mux.HandleFunc("/about", AboutHandler)
	mux.HandleFunc("/services", ServicesHandler)
	mux.HandleFunc("/catalog", CatalogHandler)
	mux.HandleFunc("/notfound", NotFoundHandler)

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
}
