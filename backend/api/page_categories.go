package api

import (
	"cookbook/repo"
	"encoding/json"
	"log"
	"net/http"
)

// рендеринг страницы категорий
func PageCategories(w http.ResponseWriter, r *http.Request) {
	log.Println("Рендеринг страницы категорий:", r.URL.Path)
	if r.URL.Path != "/categories" {
		http.NotFound(w, r)
		return
	}

	// Указываем, что возвращаем HTML
	w.Header().Set("Content-Type", "text/html")
	http.ServeFile(w, r, "../frontend/categories.html")
}

// вывод всех категорий
func PageCategoriesAll(w http.ResponseWriter, r *http.Request) {
	log.Println("Запрос к категориям:", r.URL.Path)
	if r.URL.Path != "/categories/all" || r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	// Получение последних 20 рецептов из базы данных
	categories, err := repo.GetCategory()
	if err != nil {
		http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
		log.Printf("Ошибка базы данных при получении категорий")
		return
	}

	// Указываем, что возвращаем JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}
