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
