package api

import (
	"cookbook/internal"
	"cookbook/repo"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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
	log.Println("Запрос на сохранение рецепта:", r.URL.Path)
	if r.URL.Path != "/edit/save" {
		http.NotFound(w, r)
		return
	}

	var recipe repo.SaveRecipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		http.Error(w, "Ошибка декодирования JSON", http.StatusBadRequest)
		log.Println("ошибка декодирования json")
		return
	}

	userID := internal.GetUserIDToken(w, r)

	// Получаем ID рецепта из запроса
	recipeId := r.URL.Query().Get("recipe_id")
	if recipeId == "" {
		http.Error(w, "id рецепта не указано", http.StatusBadRequest)
		return
	}

	ID, err := strconv.Atoi(recipeId)
	if err != nil {
		log.Println("Ошибка преобразования:", err)
		return
	}

	checkUserID, err := repo.GetUserIDByRecipeID(ID)
	if err != nil {
		log.Println("Ошибка базы данных:", err)
		return
	}

	if checkUserID != userID {
		http.Error(w, "у пользователя нет доступа к этому рецепту", http.StatusBadRequest)
		return
	}

	err = repo.UpdateRecipeInBd(ID, recipe, userID)
	if err != nil {
		http.Error(w, "Ошибка сохранения рецепта", http.StatusBadRequest)
		log.Println("ошибка сохранения рецепта")
		return
	}

	// Ответ после успешного добавления комментария
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
