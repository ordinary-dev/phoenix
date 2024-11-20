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
		t.Fatal(err)
	}

	if count != 0 {
		t.Fatal("user count is not zero")
	}

	// Create the first user.
	username := "test"
	password := "test"
	user, err := CreateUser(username, password)
	if err != nil {
		t.Fatal(err)
	}

	// Check password and get the user.
	dbUser, err := GetUserIfPasswordMatches(username, password)
	if err != nil {
		t.Fatal(err)
	}
	if dbUser.ID != user.ID {
		t.Fatal("wrong user id")
	}

	// Check wrong password handling.
	if _, err := GetUserIfPasswordMatches("test", "wrong-password"); err == nil {
		t.Fatal("wrong password was accepted")
	}

	// Count users again.
	count, err = CountUsers()
	if err != nil {
		t.Fatal(err)
	}

	if count != 1 {
		t.Fatal("user count is not one")
	}

	// Create session.
	session, err := CreateSession(user.ID)
	if err != nil {
		t.Errorf("can't create session: %v", err)
	}

	// Use session token.
	authorizedUser, _, err := GetUserByToken(session.Token)
	if err != nil {
		t.Errorf("can't use session token: %v", err)
	}

	if authorizedUser.ID != user.ID {
		t.Errorf("session belongs to a different user: %d != %d", authorizedUser.ID, user.ID)
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
	if err := DeleteUser(user.ID); err != nil {
		t.Fatal(err)
	}
}
