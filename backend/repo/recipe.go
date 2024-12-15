package repo

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

// удаление рецепта по id
func DeleteRecipe(recipeID string, userID int) error {
	// SQL-запрос для удаления рецепта
	query := `DELETE FROM recipes WHERE id = $1 AND user_id = $2`

	result, err := GetDB().Exec(context.Background(), query, recipeID, userID)
	if err != nil {
		log.Printf("Ошибка при удалении рецепта: %v", err)
		return err
	}

	// Проверяем, что рецепт был удален
	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		log.Printf("Рецепт с id %s не найден", recipeID)
		return err
	}

	return nil
}

type GetRecipe struct {
	UserID       int                       `json:"user_id"`
	Name         string                    `json:"name"`
	CookTime     string                    `json:"cook_time"`
	Ingredients  []IngredientsWithQuantity `json:"ingredients"`
	Instructions string                    `json:"instructions"`
	Steps        []RecipeStep              `json:"steps"`
	Categories   []string                  `json:"categories"`
	Photo        string                    `json:"photo,omitempty"` // Base64 фото
	Public       bool                      `json:"public"`
}

// Функция для получения рецепта по ID
func GetRecipeView(recipeID int) (GetRecipe, []Comment, error) {
	// SQL-запрос для получения всех данных о рецепте
	query := `
		SELECT
			r.user_id AS recipe_user_id,
			r.name AS recipe_name,
			r.cook_time::text AS recipe_cook_time,
			r.instructions,
			COALESCE(ENCODE(r.photo, 'base64'), '') AS recipe_photo,
			CASE WHEN r.is_private THEN FALSE ELSE TRUE END AS public,
			i.name AS ingredient_name,
			ri.quantity AS ingredient_quantity,
			rs.instructions AS step_instruction,
			COALESCE(ENCODE(rs.photo, 'base64'), '') AS step_photo,
			dt.name AS category_name,
			c.username AS comment_username,
			c.comment AS comment_text
		FROM recipes r
		LEFT JOIN recipe_ingredients ri ON r.id = ri.recipe_id
		LEFT JOIN ingredients i ON ri.ingredient_id = i.id
		LEFT JOIN recipe_step rs ON r.id = rs.recipe_id
		LEFT JOIN recipe_dish_types rdt ON r.id = rdt.recipe_id
		LEFT JOIN dish_types dt ON rdt.dish_types_id = dt.id
		LEFT JOIN comments c ON r.id = c.recipe_id
		WHERE r.id = $1
	`

	// Выполняем запрос
	rows, err := GetDB().Query(context.Background(), query, recipeID)
	if err != nil {
		log.Printf("Ошибка выполнения запроса GetRecipeView: %v", err)
		return GetRecipe{}, nil, err
	}
	defer rows.Close()

	// Инициализируем структуру рецепта и комментариев
	var recipe GetRecipe
	recipe.Ingredients = []IngredientsWithQuantity{}
	recipe.Steps = []RecipeStep{}
	recipe.Categories = []string{}
	comments := []Comment{}

	// Обрабатываем строки результата
	for rows.Next() {
		var (
			ingredientName     *string
			ingredientQuantity *string
			stepInstruction    *string
			stepPhoto          *string
			categoryName       *string
			commentUsername    *string
			commentText        *string
		)

		err := rows.Scan(
			&recipe.UserID,
			&recipe.Name,
			&recipe.CookTime,
			&recipe.Instructions,
			&recipe.Photo,
			&recipe.Public,
			&ingredientName,
			&ingredientQuantity,
			&stepInstruction,
			&stepPhoto,
			&categoryName,
			&commentUsername,
			&commentText,
		)
		if err != nil {
			log.Printf("Ошибка сканирования строки: %v", err)
			return GetRecipe{}, nil, err
		}

		// Добавляем ингредиенты
		if ingredientName != nil {
			ingredient := IngredientsWithQuantity{
				Name:     *ingredientName,
				Quantity: deref(ingredientQuantity),
			}
			recipe.Ingredients = appendIfUniqueIngredient(recipe.Ingredients, ingredient)
		}

		// Добавляем шаги приготовления
		if stepInstruction != nil {
			step := RecipeStep{
				Step:  *stepInstruction,
				Photo: deref(stepPhoto),
			}
			recipe.Steps = appendIfUniqueStep(recipe.Steps, step)
		}

		// Добавляем категории
		if categoryName != nil {
			recipe.Categories = appendIfUniqueCategoryName(recipe.Categories, *categoryName)
		}

		// Добавляем комментарии
		if commentUsername != nil && commentText != nil {
			comment := Comment{
				Username: *commentUsername,
				Text:     *commentText,
			}
			comments = appendIfUniqueComment(comments, comment)
		}
	}

	// Проверяем ошибки после обработки строк
	if err = rows.Err(); err != nil {
		log.Printf("Ошибка обработки строк результата: %v", err)
		return GetRecipe{}, nil, err
	}

	return recipe, comments, nil
}

// Добавляем функцию для уникальности категорий по названию
func appendIfUniqueCategoryName(categories []string, newCategory string) []string {
	for _, category := range categories {
		if category == newCategory {
			return categories
		}
	}
	return append(categories, newCategory)
}

// Структура для комментария
type Comment struct {
	Username string `json:"username"`
	Text     string `json:"text"`
}

// Вспомогательные функции
func deref(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func appendIfUniqueIngredient(ingredients []IngredientsWithQuantity, newIngredient IngredientsWithQuantity) []IngredientsWithQuantity {
	for _, ing := range ingredients {
		if ing.Name == newIngredient.Name {
			return ingredients
		}
	}
	return append(ingredients, newIngredient)
}

func appendIfUniqueStep(steps []RecipeStep, newStep RecipeStep) []RecipeStep {
	for _, step := range steps {
		if step.Step == newStep.Step {
			return steps
		}
	}
	return append(steps, newStep)
}

func appendIfUniqueComment(comments []Comment, newComment Comment) []Comment {
	for _, c := range comments {
		if c.Username == newComment.Username && c.Text == newComment.Text {
			return comments
		}
	}
	return append(comments, newComment)
}

// // Структура рецепта
// type Recipe struct {
// 	ID           int                       `json:"id"`
// 	UserID       int                       `json:"user_id"`
// 	Name         string                    `json:"name"`
// 	DishType     string                    `json:"dish_type,omitempty"`
// 	CookTime     string                    `json:"cook_time"`
// 	IsFavorite   bool                      `json:"is_favorite"`
// 	IsPrivate    bool                      `json:"is_private"`
// 	Instructions string                    `json:"instructions"`
// 	Ingredients  []IngredientsWithQuantity `json:"ingredients,omitempty"`
// 	Photo        []byte                    `json:"photo,omitempty"`
// }

// // Функция для получения рецепта по ID
// func GetRecipeView(recipeID int) (Recipe, error) {
// 	// SQL-запрос для получения рецепта и его ингредиентов
// 	query := `
// 		SELECT
// 			r.id AS recipe_id,
// 			r.user_id AS recipe_user_id,
// 			r.name AS recipe_name,
// 			dt.name AS dish_type_name,
// 			r.cook_time::text AS recipe_cook_time,
// 			r.is_favorite,
// 			r.is_private,
// 			r.instructions,
// 			r.photo,
// 			i.name AS ingredient_name,
// 			ri.quantity AS ingredient_quantity
// 		FROM recipes r
// 		LEFT JOIN recipe_ingredients ri ON r.id = ri.recipe_id
// 		LEFT JOIN ingredients i ON ri.ingredient_id = i.id
// 		LEFT JOIN recipe_dish_types rdt ON r.id = rdt.recipe_id
// 		LEFT JOIN dish_types dt ON rdt.dish_types_id = dt.id
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
// 	recipe.Ingredients = []IngredientsWithQuantity{}

// 	// Обрабатываем строки результата
// 	for rows.Next() {
// 		var ingredientName *string     // Указатель, так как может быть NULL
// 		var ingredientQuantity *string // Количество ингредиентов может быть NULL

// 		err := rows.Scan(
// 			&recipe.ID,
// 			&recipe.UserID,
// 			&recipe.Name,
// 			&recipe.DishType,
// 			&recipe.CookTime,
// 			&recipe.IsFavorite,
// 			&recipe.IsPrivate,
// 			&recipe.Instructions,
// 			&recipe.Photo,
// 			&ingredientName,
// 			&ingredientQuantity,
// 		)
// 		if err != nil {
// 			log.Printf("Ошибка сканирования строки: %v", err)
// 			return Recipe{}, err
// 		}

// 		// Добавляем ингредиенты, если они не NULL
// 		if ingredientName != nil {
// 			ingredient := IngredientsWithQuantity{
// 				Name: *ingredientName,
// 			}
// 			if ingredientQuantity != nil {
// 				ingredient.Quantity = *ingredientQuantity
// 			}
// 			recipe.Ingredients = append(recipe.Ingredients, ingredient)
// 		}
// 	}

// 	// Проверяем ошибки после завершения чтения строк
// 	if err = rows.Err(); err != nil {
// 		log.Printf("Ошибка обработки строк результата: %v", err)
// 		return Recipe{}, err
// 	}

// 	return recipe, nil
// }

// Структуры запроса
type SaveRecipe struct {
	Name         string                    `json:"name"`
	CookTime     string                    `json:"cook_time"`
	Ingredients  []IngredientsWithQuantity `json:"ingredients"`
	Instructions string                    `json:"instructions"`
	Steps        []RecipeStep              `json:"steps"`
	Categories   []int                     `json:"categories"`
	Photo        string                    `json:"photo,omitempty"` // Base64 фото
	Public       bool                      `json:"public"`
}

type IngredientsWithQuantity struct {
	Name     string `json:"name"`
	Quantity string `json:"quantity"`
}

type RecipeStep struct {
	Step  string `json:"step"`
	Photo string `json:"photo,omitempty"` // Base64 фото
}

func convertBase64ToBytes(base64String string) []byte {
	if base64String == "" {
		return nil
	}
	data, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		log.Println("Ошибка конвертации Base64:", err)
		return nil
	}
	return data
}

func SaveRecipeToBd(recipe SaveRecipe, userID int) error {
	log.Println("Сохранение рецепта")
	txOptions := pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead, // Пример изоляции
	}

	tx, err := GetDB().BeginTx(context.Background(), txOptions)
	if err != nil {
		log.Println("Ошибка начала транзакции:", err)
		return err
	}
	var recipeID int

	// добавление рецепта
	err = tx.QueryRow(context.Background(),
		`INSERT INTO recipes (user_id, name, cook_time, is_private, instructions, photo) 
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING id`,
		userID, recipe.Name, recipe.CookTime, !recipe.Public, recipe.Instructions, convertBase64ToBytes(recipe.Photo)).Scan(&recipeID)
	if err != nil {
		tx.Rollback(context.Background())
		log.Println("ошибка выполнения запроса на создание рецепта:", err)
		return errors.New("ошибка выполнения запроса на создание рецепта")
	}

	// добавление категорий
	for _, categoryID := range recipe.Categories {
		_, err := tx.Exec(context.Background(),
			`INSERT INTO recipe_dish_types (recipe_id, dish_types_id) VALUES ($1, $2)`,
			recipeID, categoryID)
		if err != nil {
			tx.Rollback(context.Background())
			log.Println("ошибка выполнения запроса на создание категорий:", err)
			return errors.New("ошибка выполнения запроса на создание категорий")
		}
	}

	// добавление ингредиентов
	for _, ing := range recipe.Ingredients {
		var ingredientID int
		err := tx.QueryRow(context.Background(),
			`INSERT INTO ingredients (name) 
			VALUES (LOWER($1)) 
			ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
			RETURNING id`,
			ing.Name).Scan(&ingredientID)
		if err != nil {
			tx.Rollback(context.Background())
			log.Println("ошибка выполнения запроса на создание ингредиентов:", err)
			return errors.New("ошибка выполнения запроса на создание ингредиентов")
		}

		_, err = tx.Exec(context.Background(),
			`INSERT INTO recipe_ingredients (recipe_id, ingredient_id, quantity) 
			VALUES ($1, $2, $3)`,
			recipeID, ingredientID, ing.Quantity)
		if err != nil {
			tx.Rollback(context.Background())
			log.Println("ошибка выполнения запроса на создание пар рецепт - ингредиент:", err)
			return errors.New("ошибка выполнения запроса на создание пар рецепт - ингредиент")
		}
	}

	// добавление этапов приготовления
	for _, step := range recipe.Steps {
		_, err = tx.Exec(context.Background(),
			`INSERT INTO recipe_step (recipe_id, instructions, photo) 
			VALUES ($1, $2, $3)`,
			recipeID, step.Step, convertBase64ToBytes(step.Photo))
		if err != nil {
			tx.Rollback(context.Background())
			log.Println("ошибка выполнения запроса на создание этапов:", err)
			return errors.New("ошибка выполнения запроса на создание этапов")
		}
	}

	// Завершение транзакции
	err = tx.Commit(context.Background())
	if err != nil {
		log.Println("Ошибка завершения транзакции:", err)
		return errors.New("ошибка завершения транзакции")
	}

	fmt.Println("успешное сохранение всех данных рецепта")
	return nil

}
