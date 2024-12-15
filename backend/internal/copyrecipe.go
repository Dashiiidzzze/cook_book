package internal

import (
	"cookbook/repo"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// Структура для парсинга JSON
type RequestDataCopy struct {
	RecipeID string `json:"recipeId"`
}

func CopyRecipe(w http.ResponseWriter, r *http.Request) {
	log.Println("Запрос на копирование рецепта:", r.URL.Path)
	if r.URL.Path != "/api/add-to-myrecipes" {
		http.NotFound(w, r)
		return
	}
	// Читаем JSON из тела запроса
	var requestData RequestDataCopy
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		log.Println("Ошибка при парсинге JSON:", err)
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	if requestData.RecipeID == "" {
		log.Println("RecipeID не указан")
		http.Error(w, "RecipeID не указан", http.StatusBadRequest)
		return
	}

	// Преобразуем ID из строки в число
	ID, err := strconv.Atoi(requestData.RecipeID)
	if err != nil {
		log.Println("Ошибка преобразования ID:", err)
		http.Error(w, "Неверный формат ID рецепта", http.StatusBadRequest)
		return
	}

	// Получаем данные рецепта
	recipe, _, err := repo.GetRecipeView(ID)
	if err != nil {
		log.Println("Ошибка получения данных из базы:", err)
		http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
		return
	}

	userID := GetUserIDToken(w, r)

	// проверка если это его собственный рецепт
	if recipe.UserID == userID {
		http.Error(w, "Пользователь с таким именем уже существует", http.StatusConflict)
		return
	}

	// Получаем ID категорий на основе их названий
	categoryMap, err := repo.GetCategoryIDsByNames(recipe.Categories)
	if err != nil {
		log.Println("Ошибка получения категорий:", err)
		http.Error(w, "Ошибка получения категорий", http.StatusInternalServerError)
		return
	}

	// Преобразуем GetRecipe в SaveRecipe
	saveRecipe := repo.SaveRecipe{
		Name:         recipe.Name,
		CookTime:     recipe.CookTime,
		Ingredients:  recipe.Ingredients,
		Instructions: recipe.Instructions,
		Steps:        recipe.Steps,
		Photo:        recipe.Photo,
		Public:       false,
	}

	// Заполняем ID категорий
	for _, categoryName := range recipe.Categories {
		if id, exists := categoryMap[categoryName]; exists {
			saveRecipe.Categories = append(saveRecipe.Categories, id)
		} else {
			log.Printf("Категория '%s' не найдена в базе данных", categoryName)
		}
	}

	err = repo.SaveRecipeToBd(saveRecipe, userID)
	if err != nil {
		http.Error(w, "Ошибка сохранения рецепта", http.StatusBadRequest)
		return
	}

	//w.Header().Set("Content-Type", "application/json")
	http.Redirect(w, r, "/myrecipes", http.StatusOK)
	//json.NewEncoder(w).Encode(string("Рецепт успешно скопирован"))
}
