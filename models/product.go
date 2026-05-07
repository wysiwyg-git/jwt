// models/product.go
package models

type Product struct {
	ID          int
	Name        string
	Category    string // "анаэробные", "цианоакрилатные", "эпоксидные"
	Description string
	Image       string // путь к картинке, например /static/images/loctite243.jpg
	// Можно добавить цену, артикул и т.д.
}

// Примерный каталог
var AllProducts = []Product{
	{
		ID: 1, Name: "Loctite 243",
		Category:    "анаэробные",
		Description: "Средний фиксатор резьбы. Устойчив к вибрации, маслу и температуре.",
		Image:       "/static/images/product-placeholder.jpg",
	},
	{
		ID: 2, Name: "Loctite 270",
		Category:    "анаэробные",
		Description: "Высокопрочный фиксатор для крупных резьб и шпилек.",
		Image:       "/static/images/product-placeholder.jpg",
	},
	{
		ID: 3, Name: "Loctite 406",
		Category:    "цианоакрилатные",
		Description: "Мгновенное склеивание металлов, пластиков, эластомеров.",
		Image:       "/static/images/product-placeholder.jpg",
	},
	{
		ID: 4, Name: "Loctite 401",
		Category:    "цианоакрилатные",
		Description: "Универсальный быстросхватывающийся клей низкой вязкости.",
		Image:       "/static/images/product-placeholder.jpg",
	},
	{
		ID: 5, Name: "Loctite EA 3450",
		Category:    "эпоксидные",
		Description: "Двухкомпонентный эпоксидный состав для ремонта металла.",
		Image:       "/static/images/product-placeholder.jpg",
	},
	{
		ID: 6, Name: "Loctite EA 9460",
		Category:    "эпоксидные",
		Description: "Высокопрочный эпоксидный клей для конструкционных соединений.",
		Image:       "/static/images/product-placeholder.jpg",
	},
	// Добавьте ещё продукты по желанию
}
