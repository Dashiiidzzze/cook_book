package internal

import (
	"cookbook/repo"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = []byte("your-secret-key") // Секретный ключ для подписи токена

// Структура для хранения данных пользователя
// type Credentials struct {
// 	Username string `json:"username"`
// 	Password string `json:"password"`
// }

// Структура данных в токене
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
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

// Middleware для проверки токена
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Читаем куку с токеном
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

		// Логируем успешную аутентификацию
		log.Printf("Пользователь %s авторизован", claims.Username)

		// Выполняем следующий обработчик
		next(w, r)
	}
}

// Генерация JWT токена
func GenerateToken(userID int, username string) (string, error) {
	// Создаем claims (данные, которые хранятся в токене)
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(), // Токен действует 3 дня
	}

	// Подписываем токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// Регистрация пользователя
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	err := repo.SaveUser(username, password)
	if err != nil {
		if err.Error() == "пользователь уже существует" {
			http.Error(w, "Пользователь с таким именем уже существует", http.StatusConflict)
			return
		}
		log.Printf("Ошибка при сохранении пользователя: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Пользователь успешно зарегистрирован"))
}

// Обработчик входа пользователя
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	if username == "" || password == "" {
		http.Error(w, "Требуется имя пользователя и  пароль", http.StatusBadRequest)
		return
	}

	// получени ехша пароля
	userID, passwdHash, err := repo.GetUser(username)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Неверное имя пользователя или пароль", http.StatusUnauthorized)
		} else {
			log.Printf("Ошибка при получении пользователя из базы: %v", err)
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		}
		return
	}

	// сравнение хеша и пароля
	err = bcrypt.CompareHashAndPassword([]byte(passwdHash), []byte(password))
	if err != nil {
		http.Error(w, "Неверное имя пользователя или пароль", http.StatusUnauthorized)
		return
	}

	token, err := GenerateToken(userID, username)
	if err != nil {
		http.Error(w, "Ошибка генерации токена", http.StatusInternalServerError)
		return
	}

	// Устанавливаем токен в cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(72 * time.Hour),
		HttpOnly: true,                    // Защита от доступа через JS (xss атаки)
		Secure:   false,                   // Только по HTTPS
		SameSite: http.SameSiteStrictMode, // защита от CSRF-атак
		Path:     "/",
	})

	fmt.Fprintf(w, "Login successful!")
}

// // Middleware для проверки токена
// func Authenticate(next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// Получаем токен из cookies
// 		cookie, err := r.Cookie("token")
// 		if err != nil {
// 			http.Error(w, "Unauthorized: No token", http.StatusUnauthorized)
// 			return
// 		}

// 		// Проверяем и парсим токен
// 		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
// 			return secretKey, nil
// 		})

// 		if err != nil || !token.Valid {
// 			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
// 			return
// 		}

// 		// Если токен валиден, продолжаем выполнение
// 		next(w, r)
// 	}
// }

// Защищенный маршрут
// func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintln(w, "Welcome to the protected area!")
// }
