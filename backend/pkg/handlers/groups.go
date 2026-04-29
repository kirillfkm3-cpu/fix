package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"social-network-backend/pkg/db"
	"github.com/gorilla/sessions"
	"time"
)

type CreateGroupRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type CreateGroupPostRequest struct {
	GroupID string  `json:"group_id"`
	Content string  `json:"content"`
	Image   *string `json:"image"`
}

type CreateEventRequest struct {
	GroupID     string  `json:"group_id"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	EventDate   string  `json:"event_date"`
}

type RespondEventRequest struct {
	EventID  string `json:"event_id"`
	Response string `json:"response"`
}

// CreateGroupHandler — создание группы
func CreateGroupHandler(database *sql.DB, store *sessions.CookieStore) http.HandlerFunc {
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

		var req CreateGroupRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Ошибка данных"})
			return
		}

		group := db.Group{
			Name:        req.Name,
			Description: req.Description,
			CreatorID:   userID,
		}

		err = db.CreateGroup(database, group)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Ошибка создания группы"})
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Группа создана"})
	}
}

// GetGroupsHandler — получение списка групп
func GetGroupsHandler(database *sql.DB, store *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		session, err := store.Get(r, "session-name")
		if err != nil || session.Values["user_id"] == nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"message": "Не авторизован"})
			return
		}

		groups, err := db.GetGroups(database)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Ошибка получения групп"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(groups)
	}
}

// JoinGroupHandler — запрос на вступление в группу
func JoinGroupHandler(database *sql.DB, store *sessions.CookieStore) http.HandlerFunc {
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
			GroupID string `json:"group_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Ошибка данных"})
			return
		}

		err = db.JoinGroup(database, req.GroupID, userID)
		if err == nil {
			// Создать notification для создателя группы
			var creatorID string
			database.QueryRow("SELECT creator_id FROM groups WHERE id = ?", req.GroupID).Scan(&creatorID)
			user, _ := db.GetUserByID(database, userID)
			db.CreateNotification(database, db.Notification{
				UserID:  creatorID,
				Type:    "group_join_request",
				GroupID: &req.GroupID,
				Message: user.FirstName + " " + user.LastName + " хочет вступить в вашу группу",
			})
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Запрос отправлен"})
	}
}

// AcceptGroupJoinHandler — принятие запроса на вступление (только создатель группы)
func AcceptGroupJoinHandler(database *sql.DB, store *sessions.CookieStore) http.HandlerFunc {
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
			GroupID string `json:"group_id"`
			UserID  string `json:"user_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Ошибка данных"})
			return
		}

		// Проверить, что currentUserID — создатель группы
		var creatorID string
		err = database.QueryRow("SELECT creator_id FROM groups WHERE id = ?", req.GroupID).Scan(&creatorID)
		if err != nil || creatorID != currentUserID {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"message": "Нет прав"})
			return
		}

		err = db.AcceptGroupJoin(database, req.GroupID, req.UserID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Ошибка принятия"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Пользователь принят"})
	}
}

// GetGroupPostsHandler — получение постов группы
func GetGroupPostsHandler(database *sql.DB, store *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		session, err := store.Get(r, "session-name")
		if err != nil || session.Values["user_id"] == nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"message": "Не авторизован"})
			return
		}
		userID := session.Values["user_id"].(string)

		groupID := r.URL.Query().Get("group_id")
		if groupID == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "group_id не указан"})
			return
		}

		// Проверить, что пользователь член группы
		isMember, err := db.IsGroupMember(database, groupID, userID)
		if err != nil || !isMember {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"message": "Не член группы"})
			return
		}

		posts, err := db.GetGroupPosts(database, groupID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Ошибка получения постов"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(posts)
	}
}

// CreateGroupPostHandler — создание поста в группе
func CreateGroupPostHandler(database *sql.DB, store *sessions.CookieStore) http.HandlerFunc {
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

		var req CreateGroupPostRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Ошибка данных"})
			return
		}

		// Проверить, что пользователь член группы
		isMember, err := db.IsGroupMember(database, req.GroupID, userID)
		if err != nil || !isMember {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"message": "Не член группы"})
			return
		}

		post := db.GroupPost{
			GroupID: req.GroupID,
			UserID:  userID,
			Content: req.Content,
			Image:   req.Image,
		}

		err = db.CreateGroupPost(database, post)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Ошибка создания поста"})
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Пост создан"})
	}
}

// CreateEventHandler — создание события
func CreateEventHandler(database *sql.DB, store *sessions.CookieStore) http.HandlerFunc {
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

		var req CreateEventRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Ошибка данных"})
			return
		}

		// Проверить, что пользователь член группы
		isMember, err := db.IsGroupMember(database, req.GroupID, userID)
		if err != nil || !isMember {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"message": "Не член группы"})
			return
		}

		eventDate, err := time.Parse("2006-01-02T15:04", req.EventDate)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Неверный формат даты"})
			return
		}

		event := db.Event{
			GroupID:     req.GroupID,
			CreatorID:   userID,
			Title:       req.Title,
			Description: req.Description,
			EventDate:   eventDate,
		}

		err = db.CreateEvent(database, event)
		if err == nil {
			// Создать notifications для членов группы
			members, _ := db.GetGroupMembers(database, req.GroupID)
			for _, member := range members {
				if member.ID != userID {
					db.CreateNotification(database, db.Notification{
						UserID:  member.ID,
						Type:    "event_created",
						GroupID: &req.GroupID,
						EventID: &event.ID,
						Message: "В группе создано новое событие: " + req.Title,
					})
				}
			}
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Событие создано"})
	}
}

// GetGroupEventsHandler — получение событий группы
func GetGroupEventsHandler(database *sql.DB, store *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		session, err := store.Get(r, "session-name")
		if err != nil || session.Values["user_id"] == nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"message": "Не авторизован"})
			return
		}
		userID := session.Values["user_id"].(string)

		groupID := r.URL.Query().Get("group_id")
		if groupID == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "group_id не указан"})
			return
		}

		// Проверить, что пользователь член группы
		isMember, err := db.IsGroupMember(database, groupID, userID)
		if err != nil || !isMember {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"message": "Не член группы"})
			return
		}

		events, err := db.GetGroupEvents(database, groupID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Ошибка получения событий"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(events)
	}
}

// RespondEventHandler — ответ на событие
func RespondEventHandler(database *sql.DB, store *sessions.CookieStore) http.HandlerFunc {
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

		var req RespondEventRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Ошибка данных"})
			return
		}

		// Проверить, что событие существует и пользователь член группы
		var groupID string
		err = database.QueryRow("SELECT group_id FROM events WHERE id = ?", req.EventID).Scan(&groupID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"message": "Событие не найдено"})
			return
		}

		isMember, err := db.IsGroupMember(database, groupID, userID)
		if err != nil || !isMember {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"message": "Не член группы"})
			return
		}

		err = db.RespondToEvent(database, req.EventID, userID, req.Response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Ошибка ответа"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Ответ сохранен"})
	}
}