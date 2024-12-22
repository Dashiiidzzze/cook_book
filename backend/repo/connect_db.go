package repo

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Глобальный пул соединений
var dbPool *pgxpool.Pool

func InitDB(connString string) {
	var err error
	dbPool, err = pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatalf("Не удается создать пул соединений: %v\n", err)
	}
	log.Println("Пул соединений к базе данных успешно создан")
}

// Закрытие пула
func CloseDB() {
	if dbPool != nil {
		dbPool.Close()
		log.Println("Пул соединений закрыт")
	}
}

// Получение пула соединений
func GetDB() *pgxpool.Pool {
	return dbPool
}
