package db

import (
	"database/sql"
	"time"
	"github.com/google/uuid"
)

type Post struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	Image     *string   `json:"image"`
	Privacy   string    `json:"privacy"`
	CreatedAt time.Time `json:"created_at"`
}

type Comment struct {
	ID        string    `json:"id"`
	PostID    string    `json:"post_id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	Image     *string   `json:"image"`
	CreatedAt time.Time `json:"created_at"`
}

type Message struct {
	ID         string    `json:"id"`
	SenderID   string    `json:"sender_id"`
	ReceiverID *string   `json:"receiver_id"`
	GroupID    *string   `json:"group_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

func CreatePost(db *sql.DB, post Post) error {
	if post.ID == "" {
		post.ID = uuid.New().String()
	}
	query := `
        INSERT INTO posts (id, user_id, content, image, privacy, created_at)
        VALUES (?, ?, ?, ?, ?, ?)
    `
	_, err := db.Exec(query, post.ID, post.UserID, post.Content, post.Image, post.Privacy, time.Now())
	return err
}

func GetPosts(db *sql.DB, userID string, currentUserID string) ([]Post, error) {
	// Получаем посты пользователя, видимые для currentUserID
	query := `
        SELECT p.id, p.user_id, p.content, p.image, p.privacy, p.created_at
        FROM posts p
        WHERE p.user_id = ?
        AND (
            p.privacy = 'public'
            OR (p.privacy = 'almost_private' AND EXISTS (
                SELECT 1 FROM followers f WHERE f.follower_id = ? AND f.following_id = p.user_id AND f.status = 'accepted'
            ))
            OR (p.privacy = 'private' AND ? = p.user_id)
        )
        ORDER BY p.created_at DESC
    `
	rows, err := db.Query(query, userID, currentUserID, currentUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.Image, &post.Privacy, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func GetFeedPosts(db *sql.DB, currentUserID string) ([]Post, error) {
	// Получаем посты от пользователей, на которых подписан currentUserID, плюс свои
	query := `
        SELECT p.id, p.user_id, p.content, p.image, p.privacy, p.created_at
        FROM posts p
        WHERE (
            p.user_id = ?
            OR EXISTS (
                SELECT 1 FROM followers f WHERE f.follower_id = ? AND f.following_id = p.user_id AND f.status = 'accepted'
            )
        )
        AND (
            p.privacy = 'public'
            OR (p.privacy = 'almost_private' AND EXISTS (
                SELECT 1 FROM followers f WHERE f.follower_id = ? AND f.following_id = p.user_id AND f.status = 'accepted'
            ))
            OR (p.privacy = 'private' AND ? = p.user_id)
        )
        ORDER BY p.created_at DESC
    `
	rows, err := db.Query(query, currentUserID, currentUserID, currentUserID, currentUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.Image, &post.Privacy, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func CreateComment(db *sql.DB, comment Comment) error {
	if comment.ID == "" {
		comment.ID = uuid.New().String()
	}
	query := `
        INSERT INTO comments (id, post_id, user_id, content, image, created_at)
        VALUES (?, ?, ?, ?, ?, ?)
    `
	_, err := db.Exec(query, comment.ID, comment.PostID, comment.UserID, comment.Content, comment.Image, time.Now())
	return err
}

func GetComments(db *sql.DB, postID string) ([]Comment, error) {
	query := `
        SELECT id, post_id, user_id, content, image, created_at
        FROM comments
        WHERE post_id = ?
        ORDER BY created_at ASC
    `
	rows, err := db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.Image, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func SaveMessage(db *sql.DB, message Message) error {
	if message.ID == "" {
		message.ID = uuid.New().String()
	}
	query := `
        INSERT INTO messages (id, sender_id, receiver_id, group_id, content, created_at)
        VALUES (?, ?, ?, ?, ?, ?)
    `
	_, err := db.Exec(query, message.ID, message.SenderID, message.ReceiverID, message.GroupID, message.Content, time.Now())
	return err
}

func GetPrivateMessages(db *sql.DB, user1, user2 string) ([]Message, error) {
	query := `
        SELECT id, sender_id, receiver_id, group_id, content, created_at
        FROM messages
        WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)
        ORDER BY created_at ASC
    `
	rows, err := db.Query(query, user1, user2, user2, user1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var message Message
		err := rows.Scan(&message.ID, &message.SenderID, &message.ReceiverID, &message.GroupID, &message.Content, &message.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func GetGroupMessages(db *sql.DB, groupID string) ([]Message, error) {
	query := `
        SELECT id, sender_id, receiver_id, group_id, content, created_at
        FROM messages
        WHERE group_id = ?
        ORDER BY created_at ASC
    `
	rows, err := db.Query(query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var message Message
		err := rows.Scan(&message.ID, &message.SenderID, &message.ReceiverID, &message.GroupID, &message.Content, &message.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}