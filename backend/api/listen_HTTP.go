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
	http.HandleFunc("/main/filters", PageMainFilters)
	//http.HandleFunc("/main/search", PageMainSearch)           // поиск на главной странице
	http.HandleFunc("/recipe", PageRecipe)          // страница просмотра рецепта
	http.HandleFunc("/recipe/view", PageRecipeView) // запрос рецепта
	http.HandleFunc("/recipe/add-comment", internal.PrivatAuthMiddleware(PageRecipeComment))
	http.HandleFunc("/categories", PageCategories)                               //страница категорий
	http.HandleFunc("/categories/all", PageCategoriesAll)                        // все категории
	http.HandleFunc("/categories-recipes", PageCategoriesRecipes)                // рецепты в категории
	http.HandleFunc("/categories-recipes/recipes", PageCategoriesRecipesView)    // рецепты в категории
	http.HandleFunc("/categories-recipes/filters", PageCategoriesRecipesFilters) // рецепты в категории

	// далее запросы требуют аутентификации
	http.HandleFunc("/myrecipes", internal.NoCacheMiddleware(internal.PrivatAuthMiddleware(PageMyRecipes)))                // страница с моими рецептами
	http.HandleFunc("/myrecipes/recipes", internal.NoCacheMiddleware(internal.PrivatAuthMiddleware(PageMyRecipesRecipes))) // мои рецепты
	http.HandleFunc("/myrecipes/recipes/", internal.PrivatAuthMiddleware(PageMyRecipesDeleteRecipe))                       // мои рецепты
	http.HandleFunc("/myrecipes/filter", internal.PrivatAuthMiddleware(PageMyRecipesFilters))

	http.HandleFunc("/edit", internal.PrivatAuthMiddleware(PageEdit)) // мои рецепты
	http.HandleFunc("/edit/save", internal.PrivatAuthMiddleware(PageEditSave))
	http.HandleFunc("/create", internal.NoCacheMiddleware(internal.PrivatAuthMiddleware(PageCreate))) // создание рецепта
	http.HandleFunc("/create/save", internal.NoCacheMiddleware(internal.PrivatAuthMiddleware(SaveCreate)))
	http.HandleFunc("/profile", internal.NoCacheMiddleware(internal.PrivatAuthMiddleware(PageProfile))) // профиль пользователя (POST, GET, DELETE)
	http.HandleFunc("/profile/username", internal.GetUsernameToken)
	http.HandleFunc("/profile/changepass", internal.ChangePasswordHandler)
	http.HandleFunc("/profile/logout", internal.LogoutHandler)
	http.HandleFunc("/auth", PageLogin)                         // профиль пользователя (POST, GET, DELETE)
	http.HandleFunc("/auth/register", internal.RegisterHandler) // Регистрация
	http.HandleFunc("/auth/login", internal.LoginHandler)       // Вход
	//http.HandleFunc("/api/check-token", internal.CheckToken)
	http.HandleFunc("/api/add-to-myrecipes", internal.PrivatAuthMiddleware(internal.CopyRecipe))
	http.HandleFunc("/api/search", internal.SearchRecipes)
	//http.HandleFunc("/protected", internal.Authenticate(internal.ProtectedHandler)) // Защищенный маршрут

	// Запуск сервера
	log.Println("Сервер запущен на порту " + strconv.Itoa(APIport) + " ...")
	if err := http.ListenAndServe(":"+strconv.Itoa(APIport), nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
