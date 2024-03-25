package database

import (
	"golang.org/x/crypto/bcrypt"
)

type Admin struct {
	ID       int
	Username string
	Bcrypt   string
}

func CountAdmins() (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM admins`
	if err := DB.QueryRow(query).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func CreateAdmin(username string, password string) (*Admin, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, err
	}

	query := `
        INSERT INTO admins(username, bcrypt)
        VALUES (?, ?)
        RETURNING id
    `

	var admin Admin
	admin.Username = username
	admin.Bcrypt = string(hash)

	err = DB.
		QueryRow(query, admin.Username, admin.Bcrypt).
		Scan(&admin.ID)

	if err != nil {
		return nil, err
	}

	return &admin, nil
}

func GetAdminIfPasswordMatches(username string, password string) (*Admin, error) {
	query := `
        SELECT id, username, bcrypt
        FROM admins
        WHERE username = ?
    `

	var admin Admin
	err := DB.
		QueryRow(query, username).
		Scan(&admin.ID, &admin.Username, &admin.Bcrypt)

	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(admin.Bcrypt), []byte(password))
	if err != nil {
		return nil, err
	}

	return &admin, nil
}

func DeleteAdmin(id int) error {
	query := `
        DELETE FROM admins
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
