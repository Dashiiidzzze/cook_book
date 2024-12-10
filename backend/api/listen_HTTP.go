package api

import (
	"cookbook/internal"
	"log"
	"net/http"
	"strconv"
)

func ListeningHTTP(APIport int) {
	// Сервируем статические файлы (CSS, JS, изображения)
	fs := http.FileServer(http.Dir("../frontend")) // Путь к вашей папке со статикой
	http.Handle("/", http.StripPrefix("/", fs))

	http.HandleFunc("/main", PageMain)                // главная страница
	http.HandleFunc("/main/recipes", PageMainRecipes) // запрос последних 10 рецептов
	//http.HandleFunc("/main/search", PageMainSearch)           // поиск на главной странице
	http.HandleFunc("/recipe", PageRecipe)                        // страница просмотра рецепта
	http.HandleFunc("/recipe/view", PageRecipeView)               // запрос рецепта
	http.HandleFunc("/categories", PageCategories)                //страница категорий
	http.HandleFunc("/categories/recipes", PageCategoriesRecipes) // рецепты в категории
	http.HandleFunc("/categories/all", PageCategoriesAll)         // все категории

	// далее запросы требуют аутентификации
	http.HandleFunc("/myrecipes", internal.AuthMiddleware(PageMyRecipes))                // страница с моими рецептами
	http.HandleFunc("/myrecipes/recipes", internal.AuthMiddleware(PageMyRecipesRecipes)) // мои рецепты
	http.HandleFunc("/create", internal.AuthMiddleware(PageCreate))                      // создание рецепта
	http.HandleFunc("/profile", internal.AuthMiddleware(PageProfile))                    // профиль пользователя (POST, GET, DELETE)
	http.HandleFunc("/auth", PageLogin)                                                  // профиль пользователя (POST, GET, DELETE)
	http.HandleFunc("/auth/register", internal.RegisterHandler)                          // Регистрация
	http.HandleFunc("/auth/login", internal.LoginHandler)                                // Вход
	//http.HandleFunc("/protected", internal.Authenticate(internal.ProtectedHandler)) // Защищенный маршрут

	// Запуск сервера
	log.Println("Сервер запущен на порту " + strconv.Itoa(APIport) + " ...")
	if err := http.ListenAndServe(":"+strconv.Itoa(APIport), nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
