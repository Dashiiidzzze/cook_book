package internal

import (
	"cookbook/repo"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Структура для смены пароля
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// Обработчик смены пароля
func ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("смена пароля")
	// Извлекаем токен и получаем имя пользователя через Middleware
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return
	}

	// Проверяем валидность токена
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		log.Println("Ошибка валидации токена:", err)
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return
	}

	// Проверяем срок действия токена
	if claims.ExpiresAt.Before(time.Now()) {
		log.Println("Токен истек")
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return
	}
	username := claims.Username

	// Декодируем JSON-запрос
	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	if req.OldPassword == "" || req.NewPassword == "" {
		http.Error(w, "Старый и новый пароли обязательны", http.StatusBadRequest)
		return
	}

	// Получаем текущий хеш пароля из БД
	_, currentPasswordHash, err := repo.GetUser(username)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Пользователь не найден", http.StatusNotFound)
			return
		}
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	// Проверяем старый пароль
	err = bcrypt.CompareHashAndPassword([]byte(currentPasswordHash), []byte(req.OldPassword))
	if err != nil {
		http.Error(w, "Старый пароль неверен", http.StatusUnauthorized)
		return
	}

	// Генерируем хеш для нового пароля
	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Ошибка при обработке пароля", http.StatusInternalServerError)
		return
	}

	// Обновляем пароль в БД
	err = repo.UpdatePassword(username, string(newPasswordHash))
	if err != nil {
		http.Error(w, "Не удалось обновить пароль", http.StatusInternalServerError)
		return
	}

	// Ответ пользователю
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Пароль успешно изменён"))
}

// Обработчик для выхода из профиля
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Удаляем cookie с токеном
	http.SetCookie(w, &http.Cookie{
		Name:     "token",                        // Имя cookie с токеном
		Value:    "",                             // Пустое значение
		Expires:  time.Now().Add(-1 * time.Hour), // Устанавливаем прошедшее время, чтобы cookie "истекла"
		HttpOnly: true,                           // Защищаем cookie от доступа через JS
		Secure:   false,                          // Не требуется для локальных запросов без HTTPS
		SameSite: http.SameSiteStrictMode,
		Path:     "/", // Доступно для всего сайта
	})

	// Переадресуем на страницу входа или главную
	http.Redirect(w, r, "/auth", http.StatusSeeOther)
}
