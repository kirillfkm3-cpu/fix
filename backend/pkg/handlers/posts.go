package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"social-network-backend/pkg/db"
	"github.com/gorilla/sessions"
)

type CreatePostRequest struct {
	Content string  `json:"content"`
	Image   *string `json:"image"`
	Privacy string  `json:"privacy"`
}

type CreateCommentRequest struct {
	PostID  string  `json:"post_id"`
	Content string  `json:"content"`
	Image   *string `json:"image"`
}

// CreatePostHandler — создание поста
func CreatePostHandler(database *sql.DB, store *sessions.CookieStore) http.HandlerFunc {
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

		var req CreatePostRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Ошибка данных"})
			return
		}

		post := db.Post{
			UserID:  userID,
			Content: req.Content,
			Image:   req.Image,
			Privacy: req.Privacy,
		}

		err = db.CreatePost(database, post)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Ошибка создания поста"})
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Пост создан"})
	}
}

// GetPostsHandler — получение постов пользователя
func GetPostsHandler(database *sql.DB, store *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		session, err := store.Get(r, "session-name")
		if err != nil || session.Values["user_id"] == nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"message": "Не авторизован"})
			return
		}
		currentUserID := session.Values["user_id"].(string)

		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			userID = currentUserID
		}

		posts, err := db.GetPosts(database, userID, currentUserID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Ошибка получения постов"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(posts)
	}
}

// GetFeedHandler — получение ленты постов
func GetFeedHandler(database *sql.DB, store *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		session, err := store.Get(r, "session-name")
		if err != nil || session.Values["user_id"] == nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"message": "Не авторизован"})
			return
		}
		currentUserID := session.Values["user_id"].(string)

		posts, err := db.GetFeedPosts(database, currentUserID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Ошибка получения ленты"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(posts)
	}
}

// CreateCommentHandler — создание комментария
func CreateCommentHandler(database *sql.DB, store *sessions.CookieStore) http.HandlerFunc {
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

		var req CreateCommentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Ошибка данных"})
			return
		}

		comment := db.Comment{
			PostID:  req.PostID,
			UserID:  userID,
			Content: req.Content,
			Image:   req.Image,
		}

		err = db.CreateComment(database, comment)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Ошибка создания комментария"})
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Комментарий создан"})
	}
}

// GetCommentsHandler — получение комментариев к посту
func GetCommentsHandler(database *sql.DB, store *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		session, err := store.Get(r, "session-name")
		if err != nil || session.Values["user_id"] == nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"message": "Не авторизован"})
			return
		}

		postID := r.URL.Query().Get("post_id")
		if postID == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "post_id не указан"})
			return
		}

		comments, err := db.GetComments(database, postID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Ошибка получения комментариев"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(comments)
	}
}