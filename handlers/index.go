package handlers

import (
	"company-site/templates"
	"net/http"
)

// IndexHandler – главная страница
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		renderNotFound(w, r)
		return
	}
	render(w, r, templates.IndexPage())
}
