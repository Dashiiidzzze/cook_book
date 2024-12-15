package api

import (
	"cookbook/internal"
	"cookbook/repo"
	"encoding/json"
	"fmt"
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

// отправка рецептов
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

// deleteRecipe обрабатывает запрос на удаление рецепта из БД
func PageMyRecipesDeleteRecipe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	recipeID := r.URL.Path[len("/myrecipes/recipes/"):]
	if recipeID == "" {
		http.Error(w, "Не указан ID рецепта", http.StatusBadRequest)
		return
	}

	// // Получаем ID рецепта из параметров запроса
	// recipeID := r.URL.Query().Get("id")
	// if recipeID == "" {
	// 	http.Error(w, "Не указан ID рецепта", http.StatusBadRequest)
	// 	return
	// }

	userID := internal.GetUserIDToken(w, r)

	err := repo.DeleteRecipe(recipeID, userID)
	if err != nil {
		http.Error(w, "Не удалось удалить рецепт", http.StatusBadRequest)
		return
	}

	// Успешный ответ
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Рецепт с ID %s успешно удален", recipeID)))
}

func PageMyRecipesFilters(w http.ResponseWriter, r *http.Request) {
	log.Println("Запрос к фильтрам:", r.URL.Path)
	if r.URL.Path != "/myrecipes/filter" {
		http.NotFound(w, r)
		return
	}

	userId := internal.GetUserIDToken(w, r)

	fmt.Println(userId)

	// Получение последних 20 рецептов из базы данных
	ingredients, err := repo.GetIngredients(&userId, nil)
	if err != nil {
		http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
		log.Printf("Ошибка базы данных при получении категорий")
		return
	}

	// Указываем, что возвращаем JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ingredients)
}
