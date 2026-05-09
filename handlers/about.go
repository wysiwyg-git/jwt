package handlers

import (
	"company-site/templates"
	"net/http"
)

// AboutHandler – страница «О компании»
func AboutHandler(w http.ResponseWriter, r *http.Request) {
	render(w, r, templates.AboutPage())
}
