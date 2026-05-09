package handlers

import (
	"company-site/templates"
	"net/http"
)

// ServicesHandler – страница «Услуги»
func ServicesHandler(w http.ResponseWriter, r *http.Request) {
	render(w, r, templates.ServicesPage())
}
