package repo

import (
	"context"
	"log"
)

// Структура для представления категории
type DishType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Функция для получения всех категорий из базы данных
func GetCategory() ([]DishType, error) {
	// SQL-запрос для получения всех категорий
	query := "SELECT id, name FROM dish_types ORDER BY name"

	// Выполняем запрос
	rows, err := GetDB().Query(context.Background(), query)
	if err != nil {
		log.Printf("Ошибка выполнения запроса GetAllDishTypes: %v", err)
		return nil, err
	}
	defer rows.Close()

	// Список для хранения категорий
	var dishTypes []DishType

	// Обрабатываем строки результата
	for rows.Next() {
		var dishType DishType
		if err := rows.Scan(&dishType.ID, &dishType.Name); err != nil {
			log.Printf("Ошибка чтения строки результата: %v", err)
			return nil, err
		}
		dishTypes = append(dishTypes, dishType)
	}

	// Проверяем на ошибки после завершения чтения строк
	if err = rows.Err(); err != nil {
		log.Printf("Ошибка обработки строк результата: %v", err)
		return nil, err
	}

	return dishTypes, nil
}

// type LastRecipe struct {
// 	ID          int          `json:"id"`
// 	Name        string       `json:"name"`
// 	CookTime    string       `json:"cook_time"`
// 	Photo       []byte       `json:"photo"`
// 	Ingredients []Ingredient `json:"ingredients"`
// }

// type Ingredient struct {
// 	Name string `json:"name"`
// }

// поиск последних 20 публичных рецептов и поиск ингредиентов для каждого
// func GetRecipesByCategory(categoriesID int) []LastRecipe {
// 	query := `
// 		SELECT
// 			r.id AS recipe_id,
// 			r.name AS recipe_name,
// 			r.cook_time::text AS recipe_cook_time,
// 			r.photo AS recipe_photo,
// 			i.name AS ingredient_name
// 		FROM recipes r
// 		LEFT JOIN recipe_ingredients ri ON r.id = ri.recipe_id
// 		LEFT JOIN ingredients i ON ri.ingredient_id = i.id
// 		WHERE r.is_private = false
// 		LIMIT 20;
// 	`

// 	rows, err := GetDB().Query(context.Background(), query)
// 	if err != nil {
// 		log.Printf("Error querying SearchAllRecipesWithIngredients: %v\n", err)
// 		return nil
// 	}
// 	defer rows.Close()

// 	// Используем map для группировки ингредиентов по ID рецептов
// 	recipeMap := make(map[int]*LastRecipe)
// 	var recipes []*LastRecipe

// 	for rows.Next() {
// 		var (
// 			recipeID       int
// 			recipeName     string
// 			recipeCookTime string
// 			recipePhoto    []byte
// 			ingredientName sql.NullString
// 		)

// 		if err := rows.Scan(&recipeID, &recipeName, &recipeCookTime, &recipePhoto, &ingredientName); err != nil {
// 			log.Printf("Error scanning row: %v\n", err)
// 			continue
// 		}

// 		// Проверяем, существует ли рецепт в map
// 		recipe, exists := recipeMap[recipeID]
// 		if !exists {
// 			recipe = &LastRecipe{
// 				ID:       recipeID,
// 				Name:     recipeName,
// 				CookTime: recipeCookTime,
// 				Photo:    recipePhoto,
// 			}
// 			recipeMap[recipeID] = recipe
// 			recipes = append(recipes, recipe)
// 		}

// 		// Если есть данные об ингредиенте, добавляем его
// 		if ingredientName.Valid {
// 			ingredient := Ingredient{
// 				Name: ingredientName.String,
// 			}
// 			recipe.Ingredients = append(recipe.Ingredients, ingredient)
// 		}
// 	}
// 	// Преобразуем массив указателей в массив объектов перед возвращением
// 	finalRecipes := make([]LastRecipe, 0, len(recipes))
// 	for _, r := range recipes {
// 		finalRecipes = append(finalRecipes, *r)
// 	}
// 	return finalRecipes
// }
