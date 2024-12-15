package repo

import (
	"context"
	"log"
)

// добавление комментария в бд
func AddComment(recipeID int, username string, comment string) error {
	_, err := GetDB().Exec(context.Background(), `
        INSERT INTO comments (recipe_id, username, comment)
        VALUES ($1, $2, $3)`, recipeID, "anonymous", comment)

	if err != nil {
		log.Printf("Ошибка выполнения запроса GetAllDishTypes: %v", err)
		return err
	}

	return nil
}
