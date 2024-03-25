package database

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestLinks(t *testing.T) {
	initTestDatabase(t)
	defer deleteTestDatabase(t)

	// Create the first group.
	group := Group{
		Name: "test",
	}
	if err := CreateGroup(&group); err != nil {
		t.Fatal(err)
	}

	// Create the first link.
	icon := "test/icon"
	link := Link{
		Name:    "test",
		Href:    "/test",
		GroupID: group.ID,
		Icon:    &icon,
	}
	if err := CreateLink(&link); err != nil {
		t.Fatal(err)
	}
	if link.ID == 0 {
		t.Fatal("link id is zero")
	}

	// Update link.
	link.Href = "/new-href"
	if err := UpdateLink(&link); err != nil {
		t.Fatal(err)
	}

	// Delete link.
	if err := DeleteLink(link.ID); err != nil {
		t.Fatal(err)
	}

	// Delete group.
	if err := DeleteGroup(group.ID); err != nil {
		t.Fatal(err)
	}
}
