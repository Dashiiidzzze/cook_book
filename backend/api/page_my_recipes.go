package api

import (
	"cookbook/internal"
	"cookbook/repo"
	"encoding/json"
	"log"
	"net/http"
)

// рендеринг страницы
func PageMyRecipes(w http.ResponseWriter, r *http.Request) {
	log.Println("Запрос к странице мои рецепты:", r.URL.Path)
	if r.URL.Path != "/myrecipes" {
		http.NotFound(w, r)
		return
	}

	// Указываем, что возвращаем HTML
	w.Header().Set("Content-Type", "text/html")
	http.ServeFile(w, r, "../frontend/myrecipes.html") // Путь к HTML-файлу
}

// отправка рецептов на главную
func PageMyRecipesRecipes(w http.ResponseWriter, r *http.Request) {
	log.Println("Запрос к моим рецептам:", r.URL.Path)
	if r.URL.Path != "/myrecipes/recipes" || r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	userID := internal.GetUserIDToken(w, r)

	// Фильтры
	filters := map[string]interface{}{
		"user_id": userID, // Фильтр по айди
		"limit":   20,     // Лимит в 10 рецептов
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
