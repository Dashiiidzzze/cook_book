package main

import (
	"cookbook/api"
	"cookbook/config"
	"cookbook/prom"
	"cookbook/repo"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	log.SetOutput(os.Stdout) // Логи будут отправляться в Docker logs
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Printf("Запускается модуль: %s", "логирование")
	dbip, APIport, DBport := config.ConfigRead()
	// Чтение переменных окружения

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Строка подключения к базе данных
	connString := "postgres://" + dbUser + ":" + dbPassword + "@" + dbip + ":" + strconv.Itoa(DBport) + "/" + dbName
	log.Printf("Строка подключения к базе данных: %s", connString)
	// Инициализация базы данных
	repo.InitDB(connString)
	defer repo.CloseDB() // Закрытие соединения при завершении программы

	// подключение prometheus

	prom.InitDBMetrics()

	// Периодическая проверка состояния БД
	go func() {
		for {
			prom.CheckDBState(repo.GetDB()) // Проверка состояния БД
			time.Sleep(10 * time.Second)    // Интервал проверки
		}
	}()

	// запуск прослушивания портов
	api.ListeningHTTP(APIport)
}
