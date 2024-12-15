package repo

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
)

type LastRecipe struct {
	ID          int          `json:"id"`
	Name        string       `json:"name"`
	CookTime    string       `json:"cook_time"`
	Photo       string       `json:"photo"`
	Ingredients []Ingredient `json:"ingredients"`
}

type Ingredient struct {
	Name string `json:"name"`
}

// возвращает рецепты из базы данных в зависимости от фильтров
func GetRecipesWithFilters(filters map[string]interface{}) ([]LastRecipe, error) {
	baseQuery := `
		SELECT DISTINCT
			r.id AS recipe_id,
			r.name AS recipe_name,
			r.cook_time::text AS recipe_cook_time,
			r.photo AS recipe_photo,
			i.name AS ingredient_name
		FROM recipes r
		LEFT JOIN recipe_ingredients ri ON r.id = ri.recipe_id
		LEFT JOIN ingredients i ON ri.ingredient_id = i.id
		LEFT JOIN recipe_dish_types rdt ON r.id = rdt.recipe_id
		LEFT JOIN dish_types dt ON rdt.dish_types_id = dt.id
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
		conditions = append(conditions, fmt.Sprintf("rdt.dish_types_id = $%d", paramCounter))
		args = append(args, categoryID)
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

	// Фильтр по ingredient_ids ИЗМЕНИТЬ
	if ingredientIDs, ok := filters["ingredient_ids"].([]string); ok {
		placeholders := []string{}
		for range ingredientIDs {
			placeholders = append(placeholders, fmt.Sprintf("$%d", paramCounter))
			paramCounter++
		}

		// Добавляем подзапрос для фильтрации рецептов
		conditions = append(conditions, fmt.Sprintf("r.id IN (SELECT ri.recipe_id FROM recipe_ingredients ri WHERE ri.ingredient_id IN (%s))", strings.Join(placeholders, ",")))

		// Добавляем аргументы для подзапроса
		for _, id := range ingredientIDs {
			args = append(args, id)
		}
	}

	if recipeName, ok := filters["recipe_name"]; ok {
		conditions = append(conditions, fmt.Sprintf("r.name ILIKE $%d", paramCounter))
		args = append(args, "%"+recipeName.(string)+"%")
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

		// Преобразуем photo в строку base64
		encodedPhoto := ""
		if recipePhoto != nil {
			encodedPhoto = base64.StdEncoding.EncodeToString(recipePhoto)
		}

		// Проверяем, существует ли рецепт в map
		recipe, exists := recipeMap[recipeID]
		if !exists {
			recipe = &LastRecipe{
				ID:       recipeID,
				Name:     recipeName,
				CookTime: recipeCookTime,
				Photo:    encodedPhoto,
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
