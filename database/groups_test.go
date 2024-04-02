package database

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestGroups(t *testing.T) {
	initTestDatabase(t)
	defer deleteTestDatabase(t)

	// Create the first group.
	group := Group{
		Name: "test",
	}
	if err := CreateGroup(&group); err != nil {
		t.Fatal(err)
	}
	if group.ID == 0 {
		t.Fatal("group id is zero")
	}

	// Update group.
	if err := UpdateGroup(group.ID, "new-name"); err != nil {
		t.Fatal(err)
	}

	// Read groups.
	groupList, err := GetGroupsWithLinks()
	if err != nil {
		t.Fatal(err)
	}

	if len(groupList) != 1 {
		t.Fatal("group list length is not one")
	}

	if groupList[0].Name != "new-name" {
		t.Fatal("wrong group name")
	}

	// Delete group.
	if err := DeleteGroup(group.ID); err != nil {
		t.Fatal(err)
	}
}
