package repo

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v5"
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

	// Фильтр по ingredient_ids
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

// удаление рецепта по id
func DeleteRecipe(recipeID string, userID int) error {
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

// сохранение рецепта в бд
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

// Обновление рецепта в базе данных
func UpdateRecipeInBd(recipeID int, recipe SaveRecipe, userID int) error {
	log.Println("Обновление рецепта")

	// Создание транзакции
	txOptions := pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead, // Пример изоляции
	}

	tx, err := GetDB().BeginTx(context.Background(), txOptions)
	if err != nil {
		log.Println("Ошибка начала транзакции:", err)
		return err
	}

	// Обновление данных рецепта
	_, err = tx.Exec(context.Background(),
		`UPDATE recipes 
        SET name = $1, cook_time = $2, is_private = $3, instructions = $4, photo = $5
        WHERE id = $6 AND user_id = $7`,
		recipe.Name, recipe.CookTime, !recipe.Public, recipe.Instructions, convertBase64ToBytes(recipe.Photo), recipeID, userID)
	if err != nil {
		tx.Rollback(context.Background())
		log.Println("Ошибка обновления рецепта:", err)
		return errors.New("ошибка обновления рецепта")
	}

	// Обновление категорий рецепта
	_, err = tx.Exec(context.Background(),
		`DELETE FROM recipe_dish_types WHERE recipe_id = $1`,
		recipeID)
	if err != nil {
		tx.Rollback(context.Background())
		log.Println("Ошибка удаления старых категорий:", err)
		return errors.New("ошибка удаления старых категорий")
	}

	for _, categoryID := range recipe.Categories {
		_, err := tx.Exec(context.Background(),
			`INSERT INTO recipe_dish_types (recipe_id, dish_types_id) VALUES ($1, $2)`,
			recipeID, categoryID)
		if err != nil {
			tx.Rollback(context.Background())
			log.Println("Ошибка добавления категорий:", err)
			return errors.New("ошибка добавления категорий")
		}
	}

	// Обновление ингредиентов рецепта
	// Удаление старых ингредиентов
	_, err = tx.Exec(context.Background(),
		`DELETE FROM recipe_ingredients WHERE recipe_id = $1`,
		recipeID)
	if err != nil {
		tx.Rollback(context.Background())
		log.Println("Ошибка удаления старых ингредиентов:", err)
		return errors.New("ошибка удаления старых ингредиентов")
	}

	// Добавление новых ингредиентов
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
			log.Println("Ошибка добавления ингредиентов:", err)
			return errors.New("ошибка добавления ингредиентов")
		}

		_, err = tx.Exec(context.Background(),
			`INSERT INTO recipe_ingredients (recipe_id, ingredient_id, quantity) 
            VALUES ($1, $2, $3)`,
			recipeID, ingredientID, ing.Quantity)
		if err != nil {
			tx.Rollback(context.Background())
			log.Println("Ошибка добавления ингредиента в рецепт:", err)
			return errors.New("ошибка добавления ингредиента в рецепт")
		}
	}

	// Обновление этапов приготовления
	_, err = tx.Exec(context.Background(),
		`DELETE FROM recipe_step WHERE recipe_id = $1`,
		recipeID)
	if err != nil {
		tx.Rollback(context.Background())
		log.Println("Ошибка удаления старых этапов приготовления:", err)
		return errors.New("ошибка удаления старых этапов приготовления")
	}

	// Добавление новых этапов приготовления
	for _, step := range recipe.Steps {
		_, err = tx.Exec(context.Background(),
			`INSERT INTO recipe_step (recipe_id, instructions, photo) 
            VALUES ($1, $2, $3)`,
			recipeID, step.Step, convertBase64ToBytes(step.Photo))
		if err != nil {
			tx.Rollback(context.Background())
			log.Println("Ошибка добавления этапа приготовления:", err)
			return errors.New("ошибка добавления этапа приготовления")
		}
	}

	// Завершение транзакции
	err = tx.Commit(context.Background())
	if err != nil {
		log.Println("Ошибка завершения транзакции:", err)
		return errors.New("ошибка завершения транзакции")
	}

	fmt.Println("Успешное обновление рецепта")
	return nil
}

// Функция для получения userID по id рецепта
func GetUserIDByRecipeID(recipeID int) (int, error) {
	// Подключение к базе данных (предполагается, что функция GetDB() возвращает соединение)
	conn, err := GetDB().Acquire(context.Background())
	if err != nil {
		log.Println("Ошибка подключения к базе данных:", err)
		return 0, err
	}
	defer conn.Release()

	var userID int
	// Запрос на получение user_id по id рецепта
	err = conn.QueryRow(context.Background(),
		`SELECT user_id FROM recipes WHERE id = $1`,
		recipeID).Scan(&userID)

	if err != nil {
		if err == pgx.ErrNoRows {
			log.Println("Рецепт не найден")
			return 0, errors.New("рецепт не найден")
		}
		log.Println("Ошибка выполнения запроса:", err)
		return 0, err
	}

	return userID, nil
}
