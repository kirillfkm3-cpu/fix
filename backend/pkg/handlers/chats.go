package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"social-network-backend/pkg/db"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // В продакшене проверить origin
	},
}

type ChatMessage struct {
	Type       string `json:"type"` // "private" or "group"
	ReceiverID string `json:"receiver_id,omitempty"`
	GroupID    string `json:"group_id,omitempty"`
	Content    string `json:"content"`
}

// ChatHandler — websocket для чатов
func ChatHandler(database *sql.DB, store *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "session-name")
		if err != nil || session.Values["user_id"] == nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"message": "Не авторизован"})
			return
		}
		userID := session.Values["user_id"].(string)

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Ошибка апгрейда:", err)
			return
		}
		defer conn.Close()

		for {
			var msg ChatMessage
			err := conn.ReadJSON(&msg)
			if err != nil {
				log.Println("Ошибка чтения:", err)
				break
			}

			message := db.Message{
				SenderID: userID,
				Content:  msg.Content,
			}

			if msg.Type == "private" {
				message.ReceiverID = &msg.ReceiverID
				// Проверить, что они following
				isFollower1, _ := db.IsFollower(database, userID, msg.ReceiverID)
				isFollower2, _ := db.IsFollower(database, msg.ReceiverID, userID)
				user, _ := db.GetUserByID(database, msg.ReceiverID)
				if !isFollower1 && !isFollower2 && !user.IsPublic {
					conn.WriteJSON(map[string]string{"error": "Нет доступа к чату"})
					continue
				}
			} else if msg.Type == "group" {
				message.GroupID = &msg.GroupID
				// Проверить членство в группе
				isMember, _ := db.IsGroupMember(database, msg.GroupID, userID)
				if !isMember {
					conn.WriteJSON(map[string]string{"error": "Не член группы"})
					continue
				}
			} else {
				conn.WriteJSON(map[string]string{"error": "Неверный тип"})
				continue
			}

			err = db.SaveMessage(database, message)
			if err != nil {
				log.Println("Ошибка сохранения:", err)
				conn.WriteJSON(map[string]string{"error": "Ошибка сохранения"})
				continue
			}

			// Отправить подтверждение
			conn.WriteJSON(map[string]string{"status": "sent"})
		}
	}
}

// GetMessagesHandler — получение сообщений
func GetMessagesHandler(database *sql.DB, store *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		session, err := store.Get(r, "session-name")
		if err != nil || session.Values["user_id"] == nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"message": "Не авторизован"})
			return
		}
		userID := session.Values["user_id"].(string)

		receiverID := r.URL.Query().Get("receiver_id")
		groupID := r.URL.Query().Get("group_id")

		var messages []db.Message
		if receiverID != "" {
			// Private messages
			isFollower1, _ := db.IsFollower(database, userID, receiverID)
			isFollower2, _ := db.IsFollower(database, receiverID, userID)
			user, _ := db.GetUserByID(database, receiverID)
			if !isFollower1 && !isFollower2 && !user.IsPublic {
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(map[string]string{"message": "Нет доступа"})
				return
			}
			messages, err = db.GetPrivateMessages(database, userID, receiverID)
		} else if groupID != "" {
			// Group messages
			isMember, _ := db.IsGroupMember(database, groupID, userID)
			if !isMember {
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(map[string]string{"message": "Не член группы"})
				return
			}
			messages, err = db.GetGroupMessages(database, groupID)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "receiver_id или group_id не указан"})
			return
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Ошибка получения сообщений"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(messages)
	}
}