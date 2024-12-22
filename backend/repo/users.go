package repo

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

//var salt string = "hehe, nfhtkjxrf cktlcndbz"

// сохраняет пользователя в базу данных
func SaveUser(name string, passwd string) error {
	query := `
		INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		RETURNING id
	`

	// Генерация хеша пароля с "cost" (сложность) 14
	hasher, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	passwordHash := string(hasher)

	var userID int

	err = GetDB().QueryRow(context.Background(), query, name, passwordHash).Scan(&userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return errors.New("пользователь уже существует")
		}
		return err
	}

	log.Printf("Пользователь успешно сохранен с ID: %d", userID)
	return nil
}

// получение информации о пользователе по имени
func GetUser(name string) (int, string, error) {
	query := "SELECT id, password_hash FROM users WHERE username = $1"

	var userID int
	var passwdHash string

	err := GetDB().QueryRow(context.Background(), query, name).Scan(&userID, &passwdHash)
	if err != nil {
		return 0, "", err
	}

	return userID, passwdHash, nil
}

// Обновление пароля пользователя
func UpdatePassword(username, newPasswordHash string) error {
	query := "UPDATE users SET password_hash = $1 WHERE username = $2"
	// Выполнение запроса
	result, err := GetDB().Exec(context.Background(), query, newPasswordHash, username)
	if err != nil {
		return errors.New("ошибка выполнения запроса на обновление пароля")
	}

	// Проверяем количество затронутых строк
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("пользователь не найден или пароль не изменился")
	}

	return nil
}
