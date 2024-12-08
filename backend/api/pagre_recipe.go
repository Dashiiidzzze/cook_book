package api

import (
	"cookbook/repo"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func PageRecipe(w http.ResponseWriter, r *http.Request) {
	log.Println("Загрузка страницы рецепта:", r.URL.Path)
	if r.URL.Path != "/recipe" {
		http.NotFound(w, r)
		return
	}

	// Указываем, что возвращаем HTML
	w.Header().Set("Content-Type", "text/html")
	http.ServeFile(w, r, "../frontend/recipe.html") // Путь к вашему HTML-файлу
}

// отправка рецептов на главную
func PageRecipeView(w http.ResponseWriter, r *http.Request) {
	log.Println("Запрос рецепта:", r.URL.Path)
	if r.URL.Path != "/recipe/view" || r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	recipeId := r.URL.Query().Get("recipe_id")
	if recipeId == "" {
		http.Error(w, "id рецепта не указано", http.StatusBadRequest)
		return
	}

	ID, err := strconv.Atoi(recipeId) // Преобразуем строку в int
	if err != nil {
		log.Println("Ошибка преобразования:", err)
		return
	}

	recipes, err := repo.GetRecipeView(ID)
	if err != nil {
		http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
		return
	}

	// Указываем, что возвращаем JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipes)
}
