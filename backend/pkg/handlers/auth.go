package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"social-network-backend/pkg/db"
	"golang.org/x/crypto/bcrypt"
	"github.com/gorilla/sessions"
)

// Структура для получения данных из формы входа
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginHandler — обработчик входа
func LoginHandler(database *sql.DB, store *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{"message": "Метод не поддерживается"})
			return
		}

		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Ошибка формата данных"})
			return
		}

		// 1. Ищем пользователя в базе по Email
		user, err := db.GetUserByEmail(database, req.Email)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"message": "Неверный email или пароль"})
			return
		}

		// 2. Сравниваем пароли
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"message": "Неверный email или пароль"})
			return
		}

		// 3. Создаем сессию
		session, err := store.Get(r, "session-name")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Ошибка сессии"})
			return
		}
		session.Values["user_id"] = user.ID
		session.Save(r, w)

		// 4. Успех!
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Вход выполнен успешно",
			"user_id": user.ID,
		})
	}
}

// RegisterHandler — обработчик регистрации
func RegisterHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var user db.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Ошибка данных"})
			return
		}

		// Хэшируем пароль
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Ошибка хэширования пароля"})
			return
		}
		user.Password = string(hashedPassword)

		err = db.CreateUser(database, user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Ошибка сохранения: возможно email уже занят"})
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Пользователь успешно создан"})
	}
}

// === НОВЫЙ ОБРАБОТЧИК ДЛЯ ПРОФИЛЯ ===

// GetProfileHandler отдает данные одного пользователя по его ID
func GetProfileHandler(database *sql.DB, store *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Получаем сессию
		session, err := store.Get(r, "session-name")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Ошибка сессии"})
			return
		}
		userIDInterface := session.Values["user_id"]
		if userIDInterface == nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"message": "Не авторизован"})
			return
		}
		currentUserID := userIDInterface.(string)

		// Получаем ID из параметров URL (например: /api/profile?id=uuid)
		profileUserID := r.URL.Query().Get("id")
		if profileUserID == "" {
			profileUserID = currentUserID // Если не указан, показываем свой профиль
		}

		// Вызываем функцию из pkg/db/user.go
		user, err := db.GetUserByID(database, profileUserID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"message": "Пользователь не найден"})
			return
		}

		// Проверяем приватность: если профиль приватный и текущий пользователь не владелец
		if !user.IsPublic && profileUserID != currentUserID {
			// Проверяем, является ли текущий пользователь follower'ом
			isFollower, err := db.IsFollower(database, currentUserID, profileUserID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"message": "Ошибка проверки доступа"})
				return
			}
			if !isFollower {
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(map[string]string{"message": "Профиль приватный"})
				return
			}
		}

		// Отправляем данные пользователя (Go сам превратит их в JSON благодаря тегам в структуре)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	}
}

// UpdateProfileHandler — обновление профиля
func UpdateProfileHandler(database *sql.DB, store *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{"message": "Метод не поддерживается"})
			return
		}

		session, err := store.Get(r, "session-name")
		if err != nil || session.Values["user_id"] == nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"message": "Не авторизован"})
			return
		}
		userID := session.Values["user_id"].(string)

		var req struct {
			IsPublic bool `json:"is_public"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Ошибка данных"})
			return
		}

		query := `UPDATE users SET is_public = ? WHERE id = ?`
		_, err = database.Exec(query, req.IsPublic, userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Ошибка обновления"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Профиль обновлен"})
	}
}

// LogoutHandler — обработчик выхода
func LogoutHandler(store *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		session, err := store.Get(r, "session-name")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Ошибка сессии"})
			return
		}
		session.Values["user_id"] = nil
		session.Options.MaxAge = -1
		session.Save(r, w)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Выход выполнен"})
	}
}