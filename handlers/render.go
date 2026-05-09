package handlers

import (
	"company-site/templates"
	"log"
	"net/http"

	"github.com/a-h/templ"
)

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
