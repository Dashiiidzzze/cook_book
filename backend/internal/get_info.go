package internal

import (
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

// Структура запроса с токеном
type TokenRequest struct {
	Token string `json:"token"`
}

// Структура ответа с именем пользователя
type UsernameResponse struct {
	Username string `json:"username"`
}

// получение userid из токена
func GetUserIDToken(w http.ResponseWriter, r *http.Request) int {
	// Читаем куку с токеном
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return 0
	}
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return 0
	}

	// Извлекаем userID из токена
	userID := claims.UserID

	return userID
}

// получение username из токена
func GetUsernameToken(w http.ResponseWriter, r *http.Request) {
	// Читаем куку с токеном
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return
	}
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return
	}

	response := UsernameResponse{
		Username: claims.Username,
	}

	// Указываем, что возвращаем JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// получение username из токена
func GetStringUsernameToken(w http.ResponseWriter, r *http.Request) string {
	// Читаем куку с токеном
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return ""
	}
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return ""
	}

	return claims.Username
}
