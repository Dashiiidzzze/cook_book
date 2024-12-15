package api

import (
	"cookbook/internal"
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

// Структура для объединённого ответа
type RecipeResponse struct {
	Recipe   repo.GetRecipe `json:"recipe"`
	Comments []repo.Comment `json:"comments"`
}

// Функция для отправки данных о рецепте
func PageRecipeView(w http.ResponseWriter, r *http.Request) {
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

	// Преобразуем ID из строки в число
	ID, err := strconv.Atoi(recipeId)
	if err != nil {
		log.Println("Ошибка преобразования ID:", err)
		http.Error(w, "Неверный формат ID рецепта", http.StatusBadRequest)
		return
	}

	// Получаем данные рецепта и комментариев
	recipe, comments, err := repo.GetRecipeView(ID)
	if err != nil {
		log.Println("Ошибка получения данных из базы:", err)
		http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
		return
	}

	// Формируем ответ
	response := RecipeResponse{
		Recipe:   recipe,
		Comments: comments,
	}

	// Устанавливаем заголовок ответа и отправляем JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Ошибка кодирования JSON:", err)
		http.Error(w, "Ошибка формирования ответа", http.StatusInternalServerError)
	}
}

// доделать
func PageRecipeComment(w http.ResponseWriter, r *http.Request) {
	log.Println("Добавление комментария:", r.URL.Path)
	if r.URL.Path != "/recipe/add_comment" || r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	// Чтение тела запроса
	var commentData struct {
		RecipeID int    `json:"recipe_id"`
		Comment  string `json:"comment"`
	}

	err := json.NewDecoder(r.Body).Decode(&commentData)
	if err != nil {
		http.Error(w, "Неверные данные", http.StatusBadRequest)
		return
	}

	username := internal.GetStringUsernameToken(w, r)

	err = repo.AddComment(commentData.RecipeID, username, commentData.Comment)
	if err != nil {
		http.Error(w, "Ошибка базы данных", http.StatusBadRequest)
		return
	}

	// Ответ после успешного добавления комментария
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// // Отправляем обновленные данные рецепта с комментариями
	// recipe, comments, err := repo.GetRecipeView(commentData.RecipeID)
	// if err != nil {
	// 	http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
	// 	return
	// }

	// response := struct {
	// 	Recipe   repo.GetRecipe `json:"recipe"`
	// 	Comments []repo.Comment `json:"comments"`
	// }{
	// 	Recipe:   recipe,
	// 	Comments: comments,
	// }

	// json.NewEncoder(w).Encode(response)
}
