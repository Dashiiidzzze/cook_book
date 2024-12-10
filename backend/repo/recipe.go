package repo

import (
	"context"
	"log"
)

// type Recipe struct {
// 	ID          int          `json:"id"`
// 	UserID int `json:"user_id"`
// 	Name        string       `json:"name"`
// 	DishType string `json:"dish_type,omitempty"`
// 	CookTime    string       `json:"cook_time"`
// 	IsFavorite  bool         `json:"is_favorite"`
// 	IsPrivate   bool         `json:"is_private"`
// 	Ingredients []Ingredient `json:"ingredients,omitempty"`
// 	Instructions string `json:"instructions"`
// 	Photo       []byte       `json:"photo,omitempty"`
// }

// // type Ingredient struct {
// // 	Name string `json:"name"`
// // }

// // Функция для получения всех категорий из базы данных
// func GetRecipeView(recipeID int) (Recipe, error) {
// 	// SQL-запрос для получения всех категорий
// 	query := `
// 		SELECT
// 			r.id AS recipe_id,
// 			r.user_id AS recipe_user_id,
// 			r.name AS recipe_name,
// 			d.name AS dish_type_name,
// 			r.cook_time::text AS recipe_cook_time,
// 			r.is_favorite,
// 			r.is_private,

// 			r.photo,
// 			i.name AS ingredient_name
// 		FROM recipes r
// 		LEFT JOIN recipe_ingredients ri ON r.id = ri.recipe_id
// 		LEFT JOIN ingredients i ON ri.ingredient_id = i.id
// 		WHERE r.id = $1
// 	`

// 	// Выполняем запрос
// 	rows, err := GetDB().Query(context.Background(), query, recipeID)
// 	if err != nil {
// 		log.Printf("Ошибка выполнения запроса GetRecipeView: %v", err)
// 		return Recipe{}, err
// 	}
// 	defer rows.Close()

// 	// Инициализируем структуру рецепта
// 	var recipe Recipe
// 	recipe.Ingredients = []Ingredient{}

// 	// Обрабатываем строки результата
// 	for rows.Next() {
// 		var ingredientName *string // Указатель, потому что ингредиенты могут быть NULL

// 		err := rows.Scan(
// 			&recipe.ID,
// 			&recipe.Name,
// 			&recipe.CookTime,
// 			&recipe.IsFavorite,
// 			&recipe.IsPrivate,
// 			&recipe.Photo,
// 			&recipe.CategoryID,
// 			&ingredientName,
// 		)
// 		if err != nil {
// 			log.Printf("Ошибка сканирования строки: %v", err)
// 			return Recipe{}, err
// 		}

// 		// Добавляем ингредиенты, если они не NULL
// 		if ingredientName != nil {
// 			recipe.Ingredients = append(recipe.Ingredients, Ingredient{Name: *ingredientName})
// 		}
// 	}

// 	// Проверяем на ошибки после завершения чтения строк
// 	if err = rows.Err(); err != nil {
// 		log.Printf("Ошибка обработки строк результата: %v", err)
// 		return Recipe{}, err
// 	}

// 	return recipe, nil
// }

// Структура ингредиента
type IngredientsWithQuantity struct {
	Name     string `json:"name"`
	Quantity string `json:"quantity,omitempty"`
}

// Структура рецепта
type Recipe struct {
	ID           int                       `json:"id"`
	UserID       int                       `json:"user_id"`
	Name         string                    `json:"name"`
	DishType     string                    `json:"dish_type,omitempty"`
	CookTime     string                    `json:"cook_time"`
	IsFavorite   bool                      `json:"is_favorite"`
	IsPrivate    bool                      `json:"is_private"`
	Instructions string                    `json:"instructions"`
	Ingredients  []IngredientsWithQuantity `json:"ingredients,omitempty"`
	Photo        []byte                    `json:"photo,omitempty"`
}

// Функция для получения рецепта по ID
func GetRecipeView(recipeID int) (Recipe, error) {
	// SQL-запрос для получения рецепта и его ингредиентов
	query := `
		SELECT
			r.id AS recipe_id,
			r.user_id AS recipe_user_id,
			r.name AS recipe_name,
			dt.name AS dish_type_name,
			r.cook_time::text AS recipe_cook_time,
			r.is_favorite,
			r.is_private,
			r.instructions,
			r.photo,
			i.name AS ingredient_name,
			ri.quantity AS ingredient_quantity
		FROM recipes r
		LEFT JOIN recipe_ingredients ri ON r.id = ri.recipe_id
		LEFT JOIN ingredients i ON ri.ingredient_id = i.id
		LEFT JOIN recipe_dish_types rdt ON r.id = rdt.recipe_id
		LEFT JOIN dish_types dt ON rdt.dish_types_id = dt.id
		WHERE r.id = $1
	`

	// Выполняем запрос
	rows, err := GetDB().Query(context.Background(), query, recipeID)
	if err != nil {
		log.Printf("Ошибка выполнения запроса GetRecipeView: %v", err)
		return Recipe{}, err
	}
	defer rows.Close()

	// Инициализируем структуру рецепта
	var recipe Recipe
	recipe.Ingredients = []IngredientsWithQuantity{}

	// Обрабатываем строки результата
	for rows.Next() {
		var ingredientName *string     // Указатель, так как может быть NULL
		var ingredientQuantity *string // Количество ингредиентов может быть NULL

		err := rows.Scan(
			&recipe.ID,
			&recipe.UserID,
			&recipe.Name,
			&recipe.DishType,
			&recipe.CookTime,
			&recipe.IsFavorite,
			&recipe.IsPrivate,
			&recipe.Instructions,
			&recipe.Photo,
			&ingredientName,
			&ingredientQuantity,
		)
		if err != nil {
			log.Printf("Ошибка сканирования строки: %v", err)
			return Recipe{}, err
		}

		// Добавляем ингредиенты, если они не NULL
		if ingredientName != nil {
			ingredient := IngredientsWithQuantity{
				Name: *ingredientName,
			}
			if ingredientQuantity != nil {
				ingredient.Quantity = *ingredientQuantity
			}
			recipe.Ingredients = append(recipe.Ingredients, ingredient)
		}
	}

	// Проверяем ошибки после завершения чтения строк
	if err = rows.Err(); err != nil {
		log.Printf("Ошибка обработки строк результата: %v", err)
		return Recipe{}, err
	}

	return recipe, nil
}
