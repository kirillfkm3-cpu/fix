package db

import (
	"database/sql"
	"time"
	"github.com/google/uuid"
)

type Notification struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	Type       string    `json:"type"`
	FromUserID *string   `json:"from_user_id"`
	GroupID    *string   `json:"group_id"`
	EventID    *string   `json:"event_id"`
	Message    string    `json:"message"`
	IsRead     bool      `json:"is_read"`
	CreatedAt  time.Time `json:"created_at"`
}

func CreateNotification(db *sql.DB, notification Notification) error {
	if notification.ID == "" {
		notification.ID = uuid.New().String()
	}
	query := `
        INSERT INTO notifications (id, user_id, type, from_user_id, group_id, event_id, message, is_read, created_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
    `
	_, err := db.Exec(query, notification.ID, notification.UserID, notification.Type, notification.FromUserID, notification.GroupID, notification.EventID, notification.Message, notification.IsRead, time.Now())
	return err
}

func GetNotifications(db *sql.DB, userID string) ([]Notification, error) {
	query := `
        SELECT id, user_id, type, from_user_id, group_id, event_id, message, is_read, created_at
        FROM notifications
        WHERE user_id = ?
        ORDER BY created_at DESC
    `
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		var notification Notification
		err := rows.Scan(&notification.ID, &notification.UserID, &notification.Type, &notification.FromUserID, &notification.GroupID, &notification.EventID, &notification.Message, &notification.IsRead, &notification.CreatedAt)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}
	return notifications, nil
}

func MarkNotificationRead(db *sql.DB, notificationID string) error {
	query := `UPDATE notifications SET is_read = 1 WHERE id = ?`
	_, err := db.Exec(query, notificationID)
	return err
}