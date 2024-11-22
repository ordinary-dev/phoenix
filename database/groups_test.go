package database

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestGroups(t *testing.T) {
	initTestDatabase(t)
	defer deleteTestDatabase(t)

	user, err := CreateUser("group-test", nil)
	if err != nil {
		t.Fatalf("error creating user: %v", err)
	}

	// Create the first group.
	group := Group{
		Name:     "test",
		Username: &user.Username,
	}
	if err := CreateGroup(&group); err != nil {
		t.Fatal(err)
	}
	if group.ID == 0 {
		t.Fatal("group id is zero")
	}

	_, err = GetGroup(group.ID)
	if err != nil {
		t.Errorf("can't get the group: %v", err)
	}

	// Update group.
	if err := UpdateGroup(group.ID, "new-name"); err != nil {
		t.Fatal(err)
	}

	// Read groups.
	groupList, err := GetGroupsWithLinks(&user.Username)
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
	err = TransferGroups(&user.Username, nil)
	if err != nil {
		t.Errorf("error when changing owner: %v", err)
	}

	// Delete group.
	if err := DeleteGroup(group.ID); err != nil {
		t.Fatal(err)
	}
}
