package api

import (
	"cookbook/repo"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func PageEdit(w http.ResponseWriter, r *http.Request) {
	log.Println("отрисовка страницы редактирования рецепта:", r.URL.Path)
	if r.URL.Path != "/edit" {
		http.NotFound(w, r)
		return
	}

	// Указываем, что возвращаем HTML
	w.Header().Set("Content-Type", "text/html")
	http.ServeFile(w, r, "../frontend/edit.html") // Путь к вашему HTML-файлу
}

func PageEditSave(w http.ResponseWriter, r *http.Request) {
	log.Println("Запрос рецепта:", r.URL.Path)
	if r.URL.Path != "/recipe/view" || r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	// Получаем ID рецепта из запроса
	recipeId := r.URL.Query().Get("recipe_id")
	if recipeId == "" {
		http.Error(w, "id рецепта не указано", http.StatusBadRequest)
		return
	}

	var recipe repo.SaveRecipe

	err := json.NewDecoder(r.Body).Decode(&recipe)
	if err != nil {
		http.Error(w, "Неверные данные", http.StatusBadRequest)
		return
	}

	fmt.Println(recipe)

	// Ответ после успешного добавления комментария
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
