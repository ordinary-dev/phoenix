package database

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestUsers(t *testing.T) {
	initTestDatabase(t)
	defer deleteTestDatabase(t)

	// We should have no users.
	count, err := CountUsers()
	if err != nil {
		t.Errorf("error counting users: %v", err)
	}

	if count != 0 {
		t.Error("the number of users is not zero")
	}

	// Create the first user.
	username := "test"
	password := "test"
	user, err := CreateUser(username, password)
	if err != nil {
		t.Errorf("error creating user: %v", err)
	}

	// Check password and get the user.
	dbUser, err := GetUserIfPasswordMatches(username, password)
	if err != nil {
		t.Errorf("error checking password: %v", err)
	}
	if dbUser.Username != user.Username {
		t.Error("wrong username")
	}

	// Check wrong password handling.
	if _, err := GetUserIfPasswordMatches("test", "wrong-password"); err == nil {
		t.Error("wrong password was accepted")
	}

	// Count users again.
	count, err = CountUsers()
	if err != nil {
		t.Errorf("error recounting users: %v", err)
	}

	if count != 1 {
		t.Error("user count is not one")
	}

	// Create session.
	session, err := CreateSession(user.Username)
	if err != nil {
		t.Errorf("error creating session: %v", err)
	}

	// Use session token.
	authorizedUser, _, err := GetUserByToken(session.Token)
	if err != nil {
		t.Errorf("can't use session token: %v", err)
	}

	if authorizedUser.Username != user.Username {
		t.Errorf("session belongs to a different user: %s != %s", authorizedUser.Username, user.Username)
	}

	_, _, err = GetUserByToken("wrong-token")
	if err == nil {
		t.Errorf("wrong token authorized someone")
	}

	// Delete session.
	err = DeleteSession(session.Token)
	if err != nil {
		t.Errorf("can't delete session: %v", err)
	}

	// Delete user.
	if err := DeleteUser(user.Username); err != nil {
		t.Errorf("error deleting user: %v", err)
	}
}
