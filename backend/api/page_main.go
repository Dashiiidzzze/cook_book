package api

import (
	"cookbook/repo"
	"encoding/json"
	"log"
	"net/http"
)

// рендеринг главной страницы
func PageMain(w http.ResponseWriter, r *http.Request) {
	log.Println("Запрос к главной странице:", r.URL.Path)
	if r.URL.Path != "/main" {
		http.NotFound(w, r)
		return
	}

	// Указываем, что возвращаем HTML
	w.Header().Set("Content-Type", "text/html")
	http.ServeFile(w, r, "../frontend/main.html") // Путь к HTML-файлу
}

// отправка рецептов на главную
func PageMainRecipes(w http.ResponseWriter, r *http.Request) {
	log.Println("Запрос к рецептам:", r.URL.Path)
	if r.URL.Path != "/main/recipes" || r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	// Фильтры
	filters := map[string]interface{}{
		"is_private": false, // Только общедоступные
		"limit":      20,    // Лимит в 10 рецептов
	}

	// Получение последних 20 рецептов из базы данных
	recipes, err := repo.GetRecipesWithFilters(filters)
	if err != nil {
		http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
		return
	}

	// Указываем, что возвращаем JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipes)
}

func PageMainFilters(w http.ResponseWriter, r *http.Request) {
	log.Println("Запрос к фильтрам:", r.URL.Path)
	if r.URL.Path != "/main/filters" {
		http.NotFound(w, r)
		return
	}

	// Получение последних 20 рецептов из базы данных
	ingredients, err := repo.GetIngredients(nil, nil)
	if err != nil {
		http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
		log.Printf("Ошибка базы данных при получении категорий")
		return
	}

	// Указываем, что возвращаем JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ingredients)
}
