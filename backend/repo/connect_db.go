package repo

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

// Глобальная переменная для подключения к базе данных
var db *pgx.Conn

// инициализация подключения
func InitDB(connString string) {
	var err error
	db, err = pgx.Connect(context.Background(), connString)
	if err != nil {
		log.Fatalf("Не удается подключиться к базе данных: %v\n", err)
	}
	log.Println("Подключение к базе данных прошло успешно")
}

// закрывает соединение с базой данных
func CloseDB() {
	if db != nil {
		if err := db.Close(context.Background()); err != nil {
			log.Printf("Ошибка при закрытии базы данных: %v\n", err)
		}
	}
}

// возвращает объект подключения для использования
func GetDB() *pgx.Conn {
	return db
}
