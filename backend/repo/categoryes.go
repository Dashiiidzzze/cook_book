package repo

import (
	"context"
	"log"
)

// Структура для представления категории
type DishType struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Photo string `json:"photo,omitempty"` // Base64 фото
}

// Функция для получения всех категорий из базы данных
func GetCategory() ([]DishType, error) {
	// SQL-запрос для получения всех категорий
	query := "SELECT id, name, COALESCE(ENCODE(photo, 'base64'), '') AS photo FROM dish_types ORDER BY name"

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
		if err := rows.Scan(&dishType.ID, &dishType.Name, &dishType.Photo); err != nil {
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

// получение ID категорий по их названиям
func GetCategoryIDsByNames(names []string) (map[string]int, error) {
	rows, err := GetDB().Query(context.Background(), "SELECT id, name FROM dish_types WHERE name = ANY($1)", names)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categoryMap := make(map[string]int)
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		}
		categoryMap[name] = id
	}
	return categoryMap, nil
}
