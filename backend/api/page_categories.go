package api

import (
	"cookbook/repo"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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
	http.ServeFile(w, r, "../frontend/categories.html") // Путь к вашему HTML-файлу
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

// вывод рецептов в категории
func PageCategoriesRecipes(w http.ResponseWriter, r *http.Request) {
	log.Println("Запрос к рецептам категории:", r.URL.Path)
	if r.URL.Path != "/categories/recipes" || r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}
	categoryId := r.URL.Query().Get("category_id")
	if categoryId == "" {
		http.Error(w, "Категория не указана", http.StatusBadRequest)
		return
	}

	// Фильтры
	filters := map[string]interface{}{
		"category_id": categoryId, // Фильтр по категории
		"is_private":  false,      // Только общедоступные
		"limit":       20,         // Лимит в 10 рецептов
	}

	recipes, err := repo.GetRecipesWithFilters(filters)
	if err != nil {
		http.Error(w, "Ошибка получения рецептов", http.StatusInternalServerError)
		log.Printf("Ошибка получения рецептов")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipes)
}

func PageCategoriesFilters(w http.ResponseWriter, r *http.Request) {
	log.Println("Запрос к фильтрам:", r.URL.Path)
	if r.URL.Path != "/categories/filters" {
		http.NotFound(w, r)
		return
	}
	categoryId := r.URL.Query().Get("category_id")
	if categoryId == "" {
		http.Error(w, "Категория не указана", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(categoryId)
	if err != nil {
		http.Error(w, "Категория не целое число", http.StatusBadRequest)
		return
	}

	// Получение последних 20 рецептов из базы данных
	ingredients, err := repo.GetIngredients(nil, &id)
	if err != nil {
		http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
		log.Printf("Ошибка базы данных при получении категорий")
		return
	}

	// Указываем, что возвращаем JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ingredients)
}
