package repo

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

// Ingredient структура для ответа
type IngredientFilter struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func GetIngredients(userID *int, categoryID *int) ([]IngredientFilter, error) {
	query := `
		SELECT DISTINCT i.id, i.name
		FROM ingredients i
		JOIN recipe_ingredients ri ON i.id = ri.ingredient_id
		JOIN recipes r ON ri.recipe_id = r.id
		LEFT JOIN recipe_dish_types rdt ON r.id = rdt.recipe_id
		LEFT JOIN dish_types dt ON rdt.dish_types_id = dt.id
		WHERE ($1::int IS NOT NULL AND r.user_id = $1)
			OR ($1::int IS NULL AND r.is_private = FALSE)
			AND ($2::int IS NULL OR dt.id = $2);
	`

	// Подготовка параметров
	var userIDParam, categoryIDParam sql.NullInt64

	// Проверяем, переданы ли параметры
	if userID != nil {
		userIDParam = sql.NullInt64{Int64: int64(*userID), Valid: true}
	} else {
		userIDParam = sql.NullInt64{Valid: false}
	}
	fmt.Println(userIDParam, "useridparam")

	if categoryID != nil {
		categoryIDParam = sql.NullInt64{Int64: int64(*categoryID), Valid: true}
	} else {
		categoryIDParam = sql.NullInt64{Valid: false}
	}

	// Выполняем запрос
	rows, err := GetDB().Query(context.Background(), query, userIDParam, categoryIDParam)
	if err != nil {
		log.Println("Ошибка выполнения запроса на получение ингредиентов:", err)
		return nil, err
	}
	defer rows.Close()

	// Собираем результаты
	var ingredients []IngredientFilter
	for rows.Next() {
		var ing IngredientFilter
		if err := rows.Scan(&ing.ID, &ing.Name); err != nil {
			log.Println("Ошибка при сканировании строки:", err)
			return nil, err
		}
		ingredients = append(ingredients, ing)
	}

	return ingredients, nil
}
