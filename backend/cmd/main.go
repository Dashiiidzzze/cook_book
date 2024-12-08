package main

import (
	"cookbook/api"
	"cookbook/config"
	"cookbook/repo"
	"log"
	"os"
	"strconv"
)

func main() {
	// Настраиваем вывод логов в файл
	// file, err := os.OpenFile("backend.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Fatalf("Ошибка открытия файла: %v", err)
	// }
	// defer file.Close()

	// log.SetOutput(file) // Устанавливаем вывод в файл
	// log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	// Настраиваем вывод логов в стандартный вывод
	log.SetOutput(os.Stdout) // Логи будут отправляться в Docker logs
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Println("Сервер запущен!")

	log.Println("Это информационное сообщение")
	log.Printf("Запускается модуль: %s", "логирование")
	dbip, APIport, DBport := config.ConfigRead()

	// // прослушивание базы данных
	// Строка подключения к базе данных
	connString := "postgres://dashidze:das10002@" + dbip + ":" + strconv.Itoa(DBport) + "/cookbook"
	// Инициализация базы данных
	repo.InitDB(connString)
	defer repo.CloseDB() // Закрытие соединения при завершении программы

	// запуск прослушивания портов
	api.ListeningHTTP(APIport)
	//api.ListeningHTTP(8080)

}
