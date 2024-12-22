package api

import (
	"cookbook/internal"
	"cookbook/prom"
	"log"
	"net/http"
	"strconv"
)

// Обработчик для /health, который возвращает статус 200 OK если сервер доступен
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func ListeningHTTP(APIport int) {
	prom.InitPrometheus()
	// Сервируем статические файлы (CSS, JS, изображения)
	fs := http.FileServer(http.Dir("../frontend"))
	http.Handle("/", http.StripPrefix("/", fs))

	http.Handle("/metrics", prom.PrometheusHandler()) // Метрики доступны на /metrics
	http.HandleFunc("/health", healthCheckHandler)    // эндпоинт для проверки состояния сервера

	http.Handle("/main", prom.MetricsMiddleware(http.HandlerFunc(PageMain)))

	//http.HandleFunc("/main", PageMain)                                                                                     // главная страница
	http.HandleFunc("/main/recipes", PageMainRecipes)                                                                      // запрос последних 20 рецептов
	http.HandleFunc("/main/filters", PageMainFilters)                                                                      // фильтры на главной странице
	http.HandleFunc("/recipe", PageRecipe)                                                                                 // страница просмотра рецепта
	http.HandleFunc("/recipe/view", PageRecipeView)                                                                        // запрос рецепта
	http.HandleFunc("/recipe/add-comment", internal.PrivatAuthMiddleware(PageRecipeComment))                               // добавление комментария
	http.HandleFunc("/categories", PageCategories)                                                                         // страница категорий
	http.HandleFunc("/categories/all", PageCategoriesAll)                                                                  // все категории
	http.HandleFunc("/categories-recipes", PageCategoriesRecipes)                                                          // страница рецепты в категории
	http.HandleFunc("/categories-recipes/recipes", PageCategoriesRecipesView)                                              // просмотр рецепты в категории
	http.HandleFunc("/categories-recipes/filters", PageCategoriesRecipesFilters)                                           // фильтры рецептов в категории
	http.HandleFunc("/myrecipes", internal.NoCacheMiddleware(internal.PrivatAuthMiddleware(PageMyRecipes)))                // страница с моими рецептами
	http.HandleFunc("/myrecipes/recipes", internal.NoCacheMiddleware(internal.PrivatAuthMiddleware(PageMyRecipesRecipes))) // мои рецепты
	http.HandleFunc("/myrecipes/recipes/", internal.PrivatAuthMiddleware(PageMyRecipesDeleteRecipe))                       // удаление рецепта
	http.HandleFunc("/myrecipes/filter", internal.PrivatAuthMiddleware(PageMyRecipesFilters))                              // фильтрация в мои рецепты
	http.HandleFunc("/edit", internal.PrivatAuthMiddleware(PageEdit))                                                      // страница редактирование
	http.HandleFunc("/edit/save", internal.PrivatAuthMiddleware(PageEditSave))                                             // сохранение отредактированного
	http.HandleFunc("/create", internal.NoCacheMiddleware(internal.PrivatAuthMiddleware(PageCreate)))                      // страница создание рецепта
	http.HandleFunc("/create/save", internal.NoCacheMiddleware(internal.PrivatAuthMiddleware(SaveCreate)))                 // сохранение созданного
	http.HandleFunc("/profile", internal.NoCacheMiddleware(internal.PrivatAuthMiddleware(PageProfile)))                    // страница профиль пользователя
	http.HandleFunc("/profile/username", internal.GetUsernameToken)                                                        // запрос имени пользователя
	http.HandleFunc("/profile/changepass", internal.ChangePasswordHandler)                                                 // смена пароля
	http.HandleFunc("/profile/logout", internal.LogoutHandler)                                                             // выход из профиля
	http.HandleFunc("/auth", PageLogin)                                                                                    // страница входа / регистрации
	http.HandleFunc("/auth/register", internal.RegisterHandler)                                                            // Регистрация
	http.HandleFunc("/auth/login", internal.LoginHandler)                                                                  // Вход
	http.HandleFunc("/registration", PageRegistration)
	http.HandleFunc("/api/add-to-myrecipes", internal.PrivatAuthMiddleware(internal.CopyRecipe)) // добавление рецепта в мои рецепты
	http.HandleFunc("/api/search", internal.SearchRecipes)                                       // поиск по названию

	// Запуск сервера
	log.Println("Сервер запущен на порту " + strconv.Itoa(APIport) + " ...")
	if err := http.ListenAndServe(":"+strconv.Itoa(APIport), nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
