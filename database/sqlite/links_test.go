package sqlite

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ordinary-dev/phoenix/database/entities"
)

func TestLinks(t *testing.T) {
	db := initTestDatabase(t)
	defer deleteTestDatabase(t)

	// Create the first group.
	group := entities.Group{
		Name: "test",
	}
	if err := db.CreateGroup(&group); err != nil {
		t.Fatal(err)
	}

	// Create the first link.
	icon := "test/icon"
	link := entities.Link{
		Name:    "test",
		Href:    "/test",
		GroupID: group.ID,
		Icon:    &icon,
	}
	if err := db.CreateLink(&link); err != nil {
		t.Fatal(err)
	}
	if link.ID == 0 {
		t.Fatal("link id is zero")
	}

	// Update link.
	link.Href = "/new-href"
	if err := db.UpdateLink(&link); err != nil {
		t.Fatal(err)
	}

	// Delete link.
	if err := db.DeleteLink(link.ID); err != nil {
		t.Fatal(err)
	}

	// Delete group.
	if err := db.DeleteGroup(group.ID); err != nil {
		t.Fatal(err)
	}
}
