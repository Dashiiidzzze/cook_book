package config

import (
	"encoding/json"
	"log"
	"os"
)

// Структура для хранения данных из JSON
type ConfigStruct struct {
	DatabaseIP   string `json:"database_ip"`
	APIPort      int    `json:"api_port"`
	DatabasePort int    `json:"database_port"`
}

// чтение конфигов из файла
func ConfigRead() (string, int, int) {
	file, err := os.ReadFile("../config/config.json")
	if err != nil {
		log.Fatalf("Не удалось открыть файл: %v", err)
	}

	var config ConfigStruct
	if err := json.Unmarshal(file, &config); err != nil {
		log.Fatalf("Ошибка при парсинге JSON: %v", err)
	}

	return config.DatabaseIP, config.APIPort, config.DatabasePort
}
