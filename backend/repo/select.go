package repo

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
)

type LastRecipe struct {
	ID          int          `json:"id"`
	Name        string       `json:"name"`
	CookTime    string       `json:"cook_time"`
	Photo       []byte       `json:"photo"`
	Ingredients []Ingredient `json:"ingredients"`
}

type Ingredient struct {
	Name string `json:"name"`
}

// поиск последних 20 публичных рецептов и поиск ингредиентов для каждого
func GetLastRecipes() []LastRecipe {
	query := `
		SELECT
			r.id AS recipe_id,
			r.name AS recipe_name,
			r.cook_time::text AS recipe_cook_time,
			r.photo AS recipe_photo,
			i.name AS ingredient_name
		FROM recipes r
		LEFT JOIN recipe_ingredients ri ON r.id = ri.recipe_id
		LEFT JOIN ingredients i ON ri.ingredient_id = i.id
		WHERE r.is_private = false
		LIMIT 20;
	`

	rows, err := GetDB().Query(context.Background(), query)
	if err != nil {
		log.Printf("Error querying SearchAllRecipesWithIngredients: %v\n", err)
		return nil
	}
	defer rows.Close()

	// Используем map для группировки ингредиентов по ID рецептов
	recipeMap := make(map[int]*LastRecipe)
	var recipes []*LastRecipe

	for rows.Next() {
		var (
			recipeID       int
			recipeName     string
			recipeCookTime string
			recipePhoto    []byte
			ingredientName sql.NullString
		)

		if err := rows.Scan(&recipeID, &recipeName, &recipeCookTime, &recipePhoto, &ingredientName); err != nil {
			log.Printf("Error scanning row: %v\n", err)
			continue
		}

		// Проверяем, существует ли рецепт в map
		recipe, exists := recipeMap[recipeID]
		if !exists {
			recipe = &LastRecipe{
				ID:       recipeID,
				Name:     recipeName,
				CookTime: recipeCookTime,
				Photo:    recipePhoto,
			}
			recipeMap[recipeID] = recipe
			recipes = append(recipes, recipe)
		}

		// Если есть данные об ингредиенте, добавляем его
		if ingredientName.Valid {
			ingredient := Ingredient{
				Name: ingredientName.String,
			}
			recipe.Ingredients = append(recipe.Ingredients, ingredient)
		}
	}
	// Преобразуем массив указателей в массив объектов перед возвращением
	finalRecipes := make([]LastRecipe, 0, len(recipes))
	for _, r := range recipes {
		finalRecipes = append(finalRecipes, *r)
	}
	return finalRecipes
}

// GetRecipesWithFilters возвращает рецепты из базы данных в зависимости от фильтров
func GetRecipesWithFilters(filters map[string]interface{}) ([]LastRecipe, error) {
	// Базовый SQL-запрос
	baseQuery := `
		SELECT
			r.id AS recipe_id,
			r.name AS recipe_name,
			r.cook_time::text AS recipe_cook_time,
			r.photo AS recipe_photo,
			i.name AS ingredient_name
		FROM recipes r
		LEFT JOIN recipe_ingredients ri ON r.id = ri.recipe_id
		LEFT JOIN ingredients i ON ri.ingredient_id = i.id
	`

	// Список условий
	var conditions []string
	var args []interface{}

	// Счетчик для параметров
	paramCounter := 1

	if userID, ok := filters["user_id"]; ok {
		conditions = append(conditions, fmt.Sprintf("r.user_id = $%d", paramCounter))
		args = append(args, userID)
		paramCounter++
	}

	if categoryID, ok := filters["category_id"]; ok {
		conditions = append(conditions, fmt.Sprintf("r.dish_type_id = $%d", paramCounter))
		args = append(args, categoryID)
		paramCounter++
	}

	if isFavorite, ok := filters["is_favorite"]; ok {
		conditions = append(conditions, fmt.Sprintf("r.is_favorite = $%d", paramCounter))
		args = append(args, isFavorite)
		paramCounter++
	}

	if isPrivate, ok := filters["is_private"]; ok {
		conditions = append(conditions, fmt.Sprintf("r.is_private = $%d", paramCounter))
		args = append(args, isPrivate)
		paramCounter++
	}

	if maxCookTime, ok := filters["max_cook_time"]; ok {
		conditions = append(conditions, fmt.Sprintf("r.cook_time <= $%d", paramCounter))
		args = append(args, maxCookTime)
		paramCounter++
	}

	// Добавляем условия в запрос
	if len(conditions) > 0 {
		baseQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Лимит
	if limit, ok := filters["limit"]; ok {
		baseQuery += fmt.Sprintf(" LIMIT $%d", paramCounter)
		args = append(args, limit)
	}

	// Выполнение запроса
	rows, err := GetDB().Query(context.Background(), baseQuery, args...)
	//rows, err := db.Query(baseQuery, args...)
	if err != nil {
		log.Printf("Ошибка выполнения запроса: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	// Используем map для группировки ингредиентов по ID рецептов
	recipeMap := make(map[int]*LastRecipe)
	var recipes []*LastRecipe

	// Обработка результатов
	for rows.Next() {
		var (
			recipeID       int
			recipeName     string
			recipeCookTime string
			recipePhoto    []byte
			ingredientName sql.NullString
		)

		if err := rows.Scan(&recipeID, &recipeName, &recipeCookTime, &recipePhoto, &ingredientName); err != nil {
			log.Printf("Error scanning row: %v\n", err)
			continue
		}

		// Проверяем, существует ли рецепт в map
		recipe, exists := recipeMap[recipeID]
		if !exists {
			recipe = &LastRecipe{
				ID:       recipeID,
				Name:     recipeName,
				CookTime: recipeCookTime,
				Photo:    recipePhoto,
			}
			recipeMap[recipeID] = recipe
			recipes = append(recipes, recipe)
		}

		// Если есть данные об ингредиенте, добавляем его
		if ingredientName.Valid {
			ingredient := Ingredient{
				Name: ingredientName.String,
			}
			recipe.Ingredients = append(recipe.Ingredients, ingredient)
		}
	}

	// Преобразуем массив указателей в массив объектов перед возвращением
	finalRecipes := make([]LastRecipe, 0, len(recipes))
	for _, r := range recipes {
		finalRecipes = append(finalRecipes, *r)
	}

	return finalRecipes, nil
}
