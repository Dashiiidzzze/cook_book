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
	http.HandleFunc("/myrecipes", PageMyRecipes)                  // страница с моими рецептами
	http.HandleFunc("/myrecipes/recipes", PageMyRecipes)          // мои рецепты
	http.HandleFunc("/create", PageCreate)                        // создание рецепта
	http.HandleFunc("/profile", PageProfile)                      // профиль пользователя (POST, GET, DELETE)
	http.HandleFunc("/login", PageLogin)                          // профиль пользователя (POST, GET, DELETE)

	http.HandleFunc("/register", internal.RegisterHandler)        // Регистрация
	http.HandleFunc("/login", loginHandler)                       // Вход
	http.HandleFunc("/protected", authenticate(protectedHandler)) // Защищенный маршрут

	// Запуск сервера
	log.Println("Сервер запущен на порту " + strconv.Itoa(APIport) + " ...")
	if err := http.ListenAndServe(":"+strconv.Itoa(APIport), nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
