package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"social-network-backend/pkg/db"
	"github.com/gorilla/sessions"
)

// GetFollowersHandler — получение подписчиков
func GetFollowersHandler(database *sql.DB, store *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		session, err := store.Get(r, "session-name")
		if err != nil || session.Values["user_id"] == nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"message": "Не авторизован"})
			return
		}

		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			userID = session.Values["user_id"].(string)
		}

		followers, err := db.GetFollowers(database, userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Ошибка получения подписчиков"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(followers)
	}
}

// GetFollowingHandler — получение подписок
func GetFollowingHandler(database *sql.DB, store *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		session, err := store.Get(r, "session-name")
		if err != nil || session.Values["user_id"] == nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"message": "Не авторизован"})
			return
		}

		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			userID = session.Values["user_id"].(string)
		}

		following, err := db.GetFollowing(database, userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Ошибка получения подписок"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(following)
	}
}

// FollowHandler — подписка на пользователя
func FollowHandler(database *sql.DB, store *sessions.CookieStore) http.HandlerFunc {
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
		followerID := session.Values["user_id"].(string)

		var req struct {
			FollowingID string `json:"following_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Ошибка данных"})
			return
		}

		// Проверяем, не подписан ли уже
		isFollower, err := db.IsFollower(database, followerID, req.FollowingID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Ошибка проверки"})
			return
		}
		if isFollower {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Уже подписаны"})
			return
		}

		// Проверяем, публичный ли профиль
		user, err := db.GetUserByID(database, req.FollowingID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"message": "Пользователь не найден"})
			return
		}

		if user.IsPublic {
			// Автоматически принимаем
			err = db.AcceptFollow(database, followerID, req.FollowingID)
		} else {
			// Создаем запрос
			err = db.FollowUser(database, followerID, req.FollowingID)
			if err == nil {
				// Создать notification
				follower, _ := db.GetUserByID(database, followerID)
				db.CreateNotification(database, db.Notification{
					UserID:     req.FollowingID,
					Type:       "follow_request",
					FromUserID: &followerID,
					Message:    follower.FirstName + " " + follower.LastName + " хочет подписаться на вас",
				})
			}
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Ошибка подписки"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Запрос отправлен"})
	}
}

// UnfollowHandler — отписка
func UnfollowHandler(database *sql.DB, store *sessions.CookieStore) http.HandlerFunc {
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
		followerID := session.Values["user_id"].(string)

		var req struct {
			FollowingID string `json:"following_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Ошибка данных"})
			return
		}

		err = db.UnfollowUser(database, followerID, req.FollowingID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Ошибка отписки"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Отписка выполнена"})
	}
}

// AcceptFollowHandler — принятие запроса на подписку
func AcceptFollowHandler(database *sql.DB, store *sessions.CookieStore) http.HandlerFunc {
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
		currentUserID := session.Values["user_id"].(string)

		var req struct {
			FollowerID string `json:"follower_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Ошибка данных"})
			return
		}

		err = db.AcceptFollow(database, req.FollowerID, currentUserID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Ошибка принятия"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Запрос принят"})
	}
}