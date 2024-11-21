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
	Username       string
	HashedPassword string
}

type Session struct {
	Token     string
	Username  string
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
        INSERT INTO users(username, hashed_password)
        VALUES (?, ?)
    `

	user := User{
		Username:       username,
		HashedPassword: string(hash),
	}

	_, err = DB.Exec(query, user.Username, user.HashedPassword)
	return &user, err
}

func GetUserIfPasswordMatches(username string, password string) (*User, error) {
	query := `
        SELECT username, hashed_password
        FROM users
        WHERE username = ?
    `

	var user User
	err := DB.
		QueryRow(query, username).
		Scan(&user.Username, &user.HashedPassword)

	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func DeleteUser(username string) error {
	query := `
        DELETE FROM users
        WHERE username = ?
    `

	res, err := DB.Exec(query, username)
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
func CreateSession(username string) (Session, error) {
	token := make([]byte, TokenLength)
	_, err := rand.Read(token)
	if err != nil {
		return Session{}, err
	}

	session := Session{
		Token:     base64.StdEncoding.EncodeToString(token),
		Username:  username,
		CreatedAt: time.Now(),
	}

	query := `
        INSERT INTO sessions(token, username, created_at)
        VALUES (?, ?, ?)
    `

	_, err = DB.Exec(query, session.Token, session.Username, session.CreatedAt.Unix())
	return session, err
}

// Find a user referenced by session token,
// if the session has not expired.
func GetUserByToken(token string) (User, Session, error) {
	query := `
        SELECT users.username, users.hashed_password, sessions.created_at
        FROM users
        INNER JOIN sessions
        ON sessions.username = users.username
        WHERE sessions.token = ?
        AND sessions.created_at >= ?
    `

	minCreatedAt := time.Now().Add(-TokenLifetime).Unix()

	var user User
	var session Session
	var createdAt int64
	err := DB.
		QueryRow(query, token, minCreatedAt).
		Scan(&user.Username, &user.HashedPassword, &createdAt)
	session.Token = token
	session.Username = user.Username
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
