package database

import (
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

func TestAdmins(t *testing.T) {
	initTestDatabase(t)
	defer deleteTestDatabase(t)

	// We should have no admins.
	count, err := CountAdmins()
	if err != nil {
		t.Fatal(err)
	}

	if count != 0 {
		t.Fatal("user count is not zero")
	}

	// Create the first user.
	username := "test"
	password := "test"
	admin, err := CreateAdmin(username, password)
	if err != nil {
		t.Fatal(err)
	}

	// Check password and get admin.
	dbAdmin, err := GetAdminIfPasswordMatches(username, password)
	if err != nil {
		t.Fatal(err)
	}
	if dbAdmin.ID != admin.ID {
		t.Fatal("wrong admin id")
	}

	// Check wrong password handling.
	if _, err := GetAdminIfPasswordMatches("test", "wrong-password"); err == nil {
		t.Fatal("wrong password was accepted")
	}

	// Count users again.
	count, err = CountAdmins()
	if err != nil {
		t.Fatal(err)
	}

	if count != 1 {
		t.Fatal("user count is not one")
	}

	// Delete user.
	if err := DeleteAdmin(admin.ID); err != nil {
		t.Fatal(err)
	}
}
