package internal

import (
	"cookbook/repo"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Обработчик запроса
func SearchRecipes(w http.ResponseWriter, r *http.Request) {
	// Парсинг параметров запроса
	searchText := r.URL.Query().Get("search")
	ingredients := r.URL.Query().Get("filters") // Фильтры через запятую
	category := r.URL.Query().Get("category")
	myrecipe := r.URL.Query().Get("myrecipe")

	// Логирование запроса
	log.Printf("Получен запрос: search=%s, filters=%s, category=%s", searchText, ingredients, category)

	// Разбор фильтров (если они переданы)
	filterMap := make(map[string]interface{})
	if searchText != "" {
		filterMap["recipe_name"] = searchText
	}
	if category != "0" {
		filterMap["category_id"] = category
	}
	if ingredients != "" {
		ingredientIDs := strings.Split(ingredients, ",")
		filterMap["ingredient_ids"] = ingredientIDs
	}
	if myrecipe == "1" {
		userID := GetUserIDToken(w, r)
		filterMap["user_id"] = userID
	} else {
		filterMap["is_private"] = false
	}

	fmt.Println(filterMap)

	// Вызов функции для получения данных из БД
	recipes, err := repo.GetRecipesWithFilters(filterMap)
	if err != nil {
		http.Error(w, "Ошибка при получении данных", http.StatusInternalServerError)
		log.Printf("Ошибка: %v", err)
		return
	}

	// Отправка результата клиенту
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(recipes); err != nil {
		http.Error(w, "Ошибка кодирования ответа", http.StatusInternalServerError)
	}
}
