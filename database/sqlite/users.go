package sqlite

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/ordinary-dev/phoenix/database/entities"
	"golang.org/x/crypto/bcrypt"
)

func (db *SqliteDB) CountUsers() (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM users`
	err := db.conn.QueryRow(query).Scan(&count)
	return count, err
}

// Create a new user.
// Ignores the operation if the user exists.
func (db *SqliteDB) CreateUser(username string, password *string) (*entities.User, error) {
	var hashedPassword *string
	if password != nil {
		rawHash, err := bcrypt.GenerateFromPassword([]byte(*password), 10)
		if err != nil {
			return nil, err
		}
		strHash := string(rawHash)
		hashedPassword = &strHash
	}

	query := `
        INSERT INTO users(username, hashed_password)
        VALUES (?, ?)
        ON CONFLICT DO NOTHING
    `

	user := entities.User{
		Username:       username,
		HashedPassword: hashedPassword,
	}

	_, err := db.conn.Exec(query, user.Username, user.HashedPassword)
	return &user, err
}

func (db *SqliteDB) GetUserIfPasswordMatches(username string, password string) (*entities.User, error) {
	query := `
        SELECT username, hashed_password
        FROM users
        WHERE username = ?
    `

	var user entities.User
	err := db.conn.
		QueryRow(query, username).
		Scan(&user.Username, &user.HashedPassword)

	if err != nil {
		return nil, err
	}

	if user.HashedPassword == nil {
		return nil, errors.New("password was not set")
	}

	err = bcrypt.CompareHashAndPassword([]byte(*user.HashedPassword), []byte(password))
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *SqliteDB) DeleteUser(username string) error {
	query := `
        DELETE FROM users
        WHERE username = ?
    `

	res, err := db.conn.Exec(query, username)
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
func (db *SqliteDB) CreateSession(username string) (entities.Session, error) {
	token := make([]byte, entities.TokenLength)
	_, err := rand.Read(token)
	if err != nil {
		return entities.Session{}, err
	}

	session := entities.Session{
		Token:     base64.StdEncoding.EncodeToString(token),
		Username:  username,
		CreatedAt: time.Now(),
	}

	query := `
        INSERT INTO sessions(token, username, created_at)
        VALUES (?, ?, ?)
    `

	_, err = db.conn.Exec(query, session.Token, session.Username, session.CreatedAt.Unix())
	return session, err
}

// Find a user referenced by session token,
// if the session has not expired.
func (db *SqliteDB) GetUserByToken(token string) (entities.User, entities.Session, error) {
	query := `
        SELECT users.username, users.hashed_password, sessions.created_at
        FROM users
        INNER JOIN sessions
        ON sessions.username = users.username
        WHERE sessions.token = ?
        AND sessions.created_at >= ?
    `

	minCreatedAt := time.Now().Add(-entities.TokenLifetime).Unix()

	var user entities.User
	var session entities.Session
	var createdAt int64
	err := db.conn.
		QueryRow(query, token, minCreatedAt).
		Scan(&user.Username, &user.HashedPassword, &createdAt)
	session.Token = token
	session.Username = user.Username
	session.CreatedAt = time.Unix(createdAt, 0)
	return user, session, err
}

// Delete user session.
func (db *SqliteDB) DeleteSession(token string) error {
	query := `
        DELETE FROM sessions
        WHERE token = ?
    `

	res, err := db.conn.Exec(query, token)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected != 1 {
		return ErrWrongNumberOfAffectedRows
	}

	return nil
}
