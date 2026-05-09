package handlers

import "net/http"

// NotFoundHandler – обработчик для /notfound
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	renderNotFound(w, r)
}
