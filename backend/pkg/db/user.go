package db

import (
	"database/sql"
	"time"
	"github.com/google/uuid"
)

type User struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	DateOfBirth string    `json:"date_of_birth"`
	Avatar      *string   `json:"avatar"`
	Nickname    *string   `json:"nickname"`
	AboutMe     *string   `json:"about_me"`
	IsPublic    bool      `json:"is_public"`
	CreatedAt   time.Time `json:"created_at"`
}

func CreateUser(db *sql.DB, user User) error {
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	query := `
        INSERT INTO users (id, email, password, first_name, last_name, date_of_birth, avatar, nickname, about_me, is_public, created_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `
	_, err := db.Exec(query,
		user.ID, user.Email, user.Password, user.FirstName, user.LastName,
		user.DateOfBirth, user.Avatar, user.Nickname, user.AboutMe, user.IsPublic,
		time.Now(),
	)
	return err
}

func GetUserByEmail(db *sql.DB, email string) (*User, error) {
	var user User
	query := `SELECT id, email, password, first_name, last_name, date_of_birth, avatar, nickname, about_me, is_public, created_at FROM users WHERE email = ?`
	err := db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName,
		&user.DateOfBirth, &user.Avatar, &user.Nickname, &user.AboutMe, &user.IsPublic, &user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByID(db *sql.DB, id string) (*User, error) {
	var user User
	query := `SELECT id, email, first_name, last_name, date_of_birth, avatar, nickname, about_me, is_public, created_at FROM users WHERE id = ?`
	
	err := db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.FirstName, &user.LastName,
		&user.DateOfBirth, &user.Avatar, &user.Nickname, &user.AboutMe, &user.IsPublic, &user.CreatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func IsFollower(db *sql.DB, followerID, followingID string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM followers WHERE follower_id = ? AND following_id = ? AND status = 'accepted'`
	err := db.QueryRow(query, followerID, followingID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetFollowers(db *sql.DB, userID string) ([]User, error) {
	query := `
        SELECT u.id, u.email, u.first_name, u.last_name, u.date_of_birth, u.avatar, u.nickname, u.about_me, u.is_public, u.created_at
        FROM users u
        JOIN followers f ON u.id = f.follower_id
        WHERE f.following_id = ? AND f.status = 'accepted'
    `
	rows, err := db.Query(query, userID)
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

func GetFollowing(db *sql.DB, userID string) ([]User, error) {
	query := `
        SELECT u.id, u.email, u.first_name, u.last_name, u.date_of_birth, u.avatar, u.nickname, u.about_me, u.is_public, u.created_at
        FROM users u
        JOIN followers f ON u.id = f.following_id
        WHERE f.follower_id = ? AND f.status = 'accepted'
    `
	rows, err := db.Query(query, userID)
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

func FollowUser(db *sql.DB, followerID, followingID string) error {
	query := `
        INSERT INTO followers (follower_id, following_id, status)
        VALUES (?, ?, 'pending')
    `
	_, err := db.Exec(query, followerID, followingID)
	return err
}

func AcceptFollow(db *sql.DB, followerID, followingID string) error {
	query := `UPDATE followers SET status = 'accepted' WHERE follower_id = ? AND following_id = ?`
	_, err := db.Exec(query, followerID, followingID)
	return err
}

func UnfollowUser(db *sql.DB, followerID, followingID string) error {
	query := `DELETE FROM followers WHERE follower_id = ? AND following_id = ?`
	_, err := db.Exec(query, followerID, followingID)
	return err
}