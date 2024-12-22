package api

import (
	"cookbook/internal"
	"cookbook/repo"
	"encoding/json"
	"log"
	"net/http"
)

// рендеринг страницы создания
func PageCreate(w http.ResponseWriter, r *http.Request) {
	log.Println("отрисовка страницы создания рецепта:", r.URL.Path)
	if r.URL.Path != "/create" {
		http.NotFound(w, r)
		return
	}

	// Указываем, что возвращаем HTML
	w.Header().Set("Content-Type", "text/html")
	http.ServeFile(w, r, "../frontend/create.html") // Путь к вашему HTML-файлу
}

// сохранение рецепта
func SaveCreate(w http.ResponseWriter, r *http.Request) {
	log.Println("Сохранение рецепта:", r.URL.Path)
	if r.URL.Path != "/create/save" {
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

	err := repo.SaveRecipeToBd(recipe, userID)
	if err != nil {
		http.Error(w, "Ошибка сохранения рецепта", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Рецепт успешно сохранен"}`))
}
