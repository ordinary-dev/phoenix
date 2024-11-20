package database

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	TokenLength   = 32
	TokenLifetime = time.Hour * 24 * 30
)

type User struct {
	ID       int
	Username string
	Bcrypt   string
}

type Session struct {
	Token     string
	UserID    int
	CreatedAt time.Time
}

func CountUsers() (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM users`
	err := DB.QueryRow(query).Scan(&count)
	return count, err
}

func CreateUser(username string, password string) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, err
	}

	query := `
        INSERT INTO users(username, bcrypt)
        VALUES (?, ?)
        RETURNING id
    `

	user := User{
		Username: username,
		Bcrypt:   string(hash),
	}

	err = DB.
		QueryRow(query, user.Username, user.Bcrypt).
		Scan(&user.ID)

	return &user, err
}

func GetUserIfPasswordMatches(username string, password string) (*User, error) {
	query := `
        SELECT id, username, bcrypt
        FROM users
        WHERE username = ?
    `

	var user User
	err := DB.
		QueryRow(query, username).
		Scan(&user.ID, &user.Username, &user.Bcrypt)

	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Bcrypt), []byte(password))
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func DeleteUser(id int) error {
	query := `
        DELETE FROM users
        WHERE id = ?
    `

	res, err := DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		return ErrWrongNumberOfAffectedRows
	}

	return nil
}

// Create a session for the user.
func CreateSession(userID int) (Session, error) {
	token := make([]byte, TokenLength)
	_, err := rand.Read(token)
	if err != nil {
		return Session{}, err
	}

	session := Session{
		Token:     base64.StdEncoding.EncodeToString(token),
		UserID:    userID,
		CreatedAt: time.Now(),
	}

	query := `
        INSERT INTO sessions(token, user_id, created_at)
        VALUES (?, ?, ?)
    `

	_, err = DB.Exec(query, session.Token, session.UserID, session.CreatedAt.Unix())
	return session, err
}

// Find a user referenced by session token,
// if the session has not expired.
func GetUserByToken(token string) (User, Session, error) {
	query := `
        SELECT users.id, users.username, users.bcrypt, sessions.created_at
        FROM users
        INNER JOIN sessions
        ON sessions.user_id = users.id
        WHERE sessions.token = ?
        AND sessions.created_at >= ?
    `

	minCreatedAt := time.Now().Add(-TokenLifetime).Unix()

	var user User
	var session Session
	var createdAt int64
	err := DB.
		QueryRow(query, token, minCreatedAt).
		Scan(&user.ID, &user.Username, &user.Bcrypt, &createdAt)
	session.Token = token
	session.UserID = user.ID
	session.CreatedAt = time.Unix(createdAt, 0)
	return user, session, err
}

// Delete user session.
func DeleteSession(token string) error {
	query := `
        DELETE FROM sessions
        WHERE token = ?
    `

	res, err := DB.Exec(query, token)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected != 1 {
		return ErrWrongNumberOfAffectedRows
	}

	return nil
}
