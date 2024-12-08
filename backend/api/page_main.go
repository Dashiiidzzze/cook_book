package api

import (
	"cookbook/repo"
	"encoding/json"
	"log"
	"net/http"
)

// type MainSearchRequest struct {
// 	Name      string `json:"name"`
// 	Dish_type string `json:"dish_type"`
// 	Cook_time string `json:"cook_time"`
// 	Holiday   string `json:"holiday"`
// }

// type MainSearchResponse struct {
// 	ID          int       `json:"id"`
// 	Name        string    `json:"name"`
// 	Cook_time   time.Time `json:"cook_time"`
// 	Ingredients []string  `json:"Ingredients"`
// }

type LastRecipe struct {
	ID          int          `json:"id"`
	Name        string       `json:"name"`
	CookTime    string       `json:"cook_time"`
	Ingredients []Ingredient `json:"ingredients"`
}

type Ingredient struct {
	Name string `json:"name"`
}

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

	// Получение последних 20 рецептов из базы данных
	recipes := repo.GetLastRecipes()
	if recipes == nil {
		http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
		return
	}

	// Указываем, что возвращаем JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipes)
}

// // поиск на главной странице
// func PageMainSearch(w http.ResponseWriter, r *http.Request) {
// 	// Парсинг JSON-запроса
// 	var req MainSearchRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, "Invalid JSON", http.StatusBadRequest)
// 		//log.Println("Это информационное сообщение")
// 		return
// 	}

// 	// recipes := internal.Search(MainSearchRequest.Name, MainSearchRequest.Dish_type, MainSearchRequest.Holiday, MainSearchRequest.Cook_time)

// 	// w.Header().Set("Content-Type", "application/json")
// 	// json.NewEncoder(w).Encode(recipes)
// }
