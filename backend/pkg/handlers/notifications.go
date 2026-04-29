package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"social-network-backend/pkg/db"
	"github.com/gorilla/sessions"
)

// GetNotificationsHandler — получение уведомлений
func GetNotificationsHandler(database *sql.DB, store *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		session, err := store.Get(r, "session-name")
		if err != nil || session.Values["user_id"] == nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"message": "Не авторизован"})
			return
		}
		userID := session.Values["user_id"].(string)

		notifications, err := db.GetNotifications(database, userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Ошибка получения уведомлений"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(notifications)
	}
}

// MarkNotificationReadHandler — отметить как прочитанное
func MarkNotificationReadHandler(database *sql.DB, store *sessions.CookieStore) http.HandlerFunc {
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

		var req struct {
			NotificationID string `json:"notification_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Ошибка данных"})
			return
		}

		err = db.MarkNotificationRead(database, req.NotificationID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Ошибка"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Отмечено как прочитанное"})
	}
}