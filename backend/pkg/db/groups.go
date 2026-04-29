package db

import (
	"database/sql"
	"time"
	"github.com/google/uuid"
)

type Group struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	CreatorID   string    `json:"creator_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type GroupPost struct {
	ID        string    `json:"id"`
	GroupID   string    `json:"group_id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	Image     *string   `json:"image"`
	CreatedAt time.Time `json:"created_at"`
}

type Event struct {
	ID          string    `json:"id"`
	GroupID     string    `json:"group_id"`
	CreatorID   string    `json:"creator_id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	EventDate   time.Time `json:"event_date"`
	CreatedAt   time.Time `json:"created_at"`
}

func CreateGroup(db *sql.DB, group Group) error {
	if group.ID == "" {
		group.ID = uuid.New().String()
	}
	query := `
        INSERT INTO groups (id, name, description, creator_id, created_at)
        VALUES (?, ?, ?, ?, ?)
    `
	_, err := db.Exec(query, group.ID, group.Name, group.Description, group.CreatorID, time.Now())
	return err
}

func GetGroups(db *sql.DB) ([]Group, error) {
	query := `SELECT id, name, description, creator_id, created_at FROM groups ORDER BY created_at DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []Group
	for rows.Next() {
		var group Group
		err := rows.Scan(&group.ID, &group.Name, &group.Description, &group.CreatorID, &group.CreatedAt)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}
	return groups, nil
}

func JoinGroup(db *sql.DB, groupID, userID string) error {
	query := `INSERT INTO group_members (group_id, user_id, status) VALUES (?, ?, 'pending')`
	_, err := db.Exec(query, groupID, userID)
	return err
}

func AcceptGroupJoin(db *sql.DB, groupID, userID string) error {
	query := `UPDATE group_members SET status = 'accepted' WHERE group_id = ? AND user_id = ?`
	_, err := db.Exec(query, groupID, userID)
	return err
}

func IsGroupMember(db *sql.DB, groupID, userID string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM group_members WHERE group_id = ? AND user_id = ? AND status = 'accepted'`
	err := db.QueryRow(query, groupID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetGroupMembers(db *sql.DB, groupID string) ([]User, error) {
	query := `
        SELECT u.id, u.email, u.first_name, u.last_name, u.date_of_birth, u.avatar, u.nickname, u.about_me, u.is_public, u.created_at
        FROM users u
        JOIN group_members gm ON u.id = gm.user_id
        WHERE gm.group_id = ? AND gm.status = 'accepted'
    `
	rows, err := db.Query(query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.DateOfBirth, &user.Avatar, &user.Nickname, &user.AboutMe, &user.IsPublic, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func CreateGroupPost(db *sql.DB, post GroupPost) error {
	if post.ID == "" {
		post.ID = uuid.New().String()
	}
	query := `
        INSERT INTO group_posts (id, group_id, user_id, content, image, created_at)
        VALUES (?, ?, ?, ?, ?, ?)
    `
	_, err := db.Exec(query, post.ID, post.GroupID, post.UserID, post.Content, post.Image, time.Now())
	return err
}

func GetGroupPosts(db *sql.DB, groupID string) ([]GroupPost, error) {
	query := `SELECT id, group_id, user_id, content, image, created_at FROM group_posts WHERE group_id = ? ORDER BY created_at DESC`
	rows, err := db.Query(query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []GroupPost
	for rows.Next() {
		var post GroupPost
		err := rows.Scan(&post.ID, &post.GroupID, &post.UserID, &post.Content, &post.Image, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func CreateEvent(db *sql.DB, event Event) error {
	if event.ID == "" {
		event.ID = uuid.New().String()
	}
	query := `
        INSERT INTO events (id, group_id, creator_id, title, description, event_date, created_at)
        VALUES (?, ?, ?, ?, ?, ?, ?)
    `
	_, err := db.Exec(query, event.ID, event.GroupID, event.CreatorID, event.Title, event.Description, event.EventDate, time.Now())
	return err
}

func GetGroupEvents(db *sql.DB, groupID string) ([]Event, error) {
	query := `SELECT id, group_id, creator_id, title, description, event_date, created_at FROM events WHERE group_id = ? ORDER BY event_date ASC`
	rows, err := db.Query(query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.GroupID, &event.CreatorID, &event.Title, &event.Description, &event.EventDate, &event.CreatedAt)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func RespondToEvent(db *sql.DB, eventID, userID, response string) error {
	query := `
        INSERT INTO event_responses (event_id, user_id, response, responded_at)
        VALUES (?, ?, ?, ?)
        ON CONFLICT(event_id, user_id) DO UPDATE SET response = excluded.response, responded_at = excluded.responded_at
    `
	_, err := db.Exec(query, eventID, userID, response, time.Now())
	return err
}