package handlers

import (
	"company-site/config"
	"company-site/mailer"
	"company-site/models"
	"company-site/templates"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/a-h/templ"
)

var (
	submissionsMu sync.Mutex
	submissions   = make(map[string]time.Time)
)

const productsPerPage = 6

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

// IndexHandler – главная страница
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		renderNotFound(w, r)
		return
	}
	render(w, r, templates.IndexPage())
}

// AboutHandler – страница «О компании»
func AboutHandler(w http.ResponseWriter, r *http.Request) {
	render(w, r, templates.AboutPage())
}

// ServicesHandler – страница «Услуги»
func ServicesHandler(w http.ResponseWriter, r *http.Request) {
	render(w, r, templates.ServicesPage())
}

// ContactHandler – страница с контактной формой
func ContactHandler(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	switch r.Method {
	case http.MethodGet:
		success := r.URL.Query().Get("success")
		var data models.ContactFormData
		render(w, r, templates.ContactPage(data, success))
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// rate limiting: не чаще раза в минуту с одного IP
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr // fallback
		}

		submissionsMu.Lock()
		last, exists := submissions[ip]
		submissionsMu.Unlock()
		if exists && time.Since(last) < time.Minute {
			http.Error(w, "Слишком много запросов. Пожалуйста, подождите 1 минуту.", http.StatusTooManyRequests)
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
			// Отправка письма
			body := mailer.BuildContactBody(data.Name, data.Company, data.Email, data.Phone, data.Message)
			if err := mailer.Send(cfg.SMTP, "Новая заявка с сайта ПромКлей", body); err != nil {
				log.Printf("Ошибка отправки email: %v", err)
				http.Redirect(w, r, "/contact?success=0", http.StatusSeeOther)
				return
			}

			submissionsMu.Lock()
			submissions[ip] = time.Now()
			submissionsMu.Unlock()

			// Успех
			http.Redirect(w, r, "/contact?success=1", http.StatusSeeOther)
			return
		}
		// Ошибки валидации – показываем форму снова
		render(w, r, templates.ContactPage(data, ""))
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

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

// NotFoundHandler – обработчик для /notfound
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	renderNotFound(w, r)
}
