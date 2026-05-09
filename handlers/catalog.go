package handlers

import (
	"company-site/models"
	"company-site/templates"
	"net/http"
	"strconv"
)

const productsPerPage = 6

// CatalogHandler – каталог продуктов
func CatalogHandler(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	pageStr := r.URL.Query().Get("page")
	page := 1
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}

	// Фильтрация
	var filtered []models.Product
	for _, p := range models.AllProducts {
		if category == "" || p.Category == category {
			filtered = append(filtered, p)
		}
	}

	// Пагинация
	totalProducts := len(filtered)
	totalPages := (totalProducts + productsPerPage - 1) / productsPerPage
	if totalPages < 1 {
		totalPages = 1
	}
	if page > totalPages {
		page = totalPages
	}
	start := (page - 1) * productsPerPage
	end := start + productsPerPage
	if end > totalProducts {
		end = totalProducts
	}
	pageProducts := filtered[start:end]

	render(w, r, templates.CatalogPage(pageProducts, category, page, totalPages))
}
