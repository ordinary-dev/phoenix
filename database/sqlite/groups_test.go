package sqlite

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ordinary-dev/phoenix/database/entities"
)

func TestGroups(t *testing.T) {
	db := initTestDatabase(t)
	defer deleteTestDatabase(t)

	user, err := db.CreateUser("group-test", nil)
	if err != nil {
		t.Fatalf("error creating user: %v", err)
	}

	// Create the first group.
	group := entities.Group{
		Name:     "test",
		Username: &user.Username,
	}
	if err := db.CreateGroup(&group); err != nil {
		t.Fatal(err)
	}
	if group.ID == 0 {
		t.Fatal("group id is zero")
	}

	_, err = db.GetGroup(group.ID)
	if err != nil {
		t.Errorf("can't get the group: %v", err)
	}

	// Update group.
	if err := db.UpdateGroup(group.ID, "new-name"); err != nil {
		t.Fatal(err)
	}

	// Read groups.
	groupList, err := db.GetGroupsWithLinks(&user.Username)
	if err != nil {
		t.Fatal(err)
	}

	if len(groupList) != 1 {
		t.Fatal("group list length is not one")
	}

	if groupList[0].Name != "new-name" {
		t.Fatal("wrong group name")
	}

	// Transfer ownership.
	err = db.TransferGroups(&user.Username, nil)
	if err != nil {
		t.Errorf("error when changing owner: %v", err)
	}

	// Delete group.
	if err := db.DeleteGroup(group.ID); err != nil {
		t.Fatal(err)
	}
}
