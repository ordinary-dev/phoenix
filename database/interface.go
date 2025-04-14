package database

import (
	"github.com/ordinary-dev/phoenix/config"
	"github.com/ordinary-dev/phoenix/database/entities"
)

type Database interface {
	// Establish database connection.
	// This method is called before any other operations.
	Connect(cfg *config.Config) error
	// Migrate database schema.
	// This method is called right after Connect() before any other operations.
	Migrate() error

	GroupsDB
	LinksDB
	UsersDB
}

type GroupsDB interface {
	GetGroupsWithLinks(username *string) ([]entities.Group, error)
	// This function should fill in the ID.
	CreateGroup(group *entities.Group) error
	// Get group by ID.
	// Child link list may be empty..
	GetGroup(id int) (entities.Group, error)
	TransferGroups(initialUsername, newUsername *string) error
	UpdateGroup(id int, name string) error
	DeleteGroup(groupID int) error
}

type LinksDB interface {
	GetLinksFromGroup(groupID int) ([]entities.Link, error)
	GetLink(id int) (*entities.Link, error)
	// This function should fill in the ID.
	CreateLink(link *entities.Link) error
	UpdateLink(link *entities.Link) error
	DeleteLink(linkID int) error
}

type UsersDB interface {
	CountUsers() (int64, error)
	// Should be a no-op if the user exists.
	CreateUser(username string, password *string) (*entities.User, error)
	GetUserIfPasswordMatches(username string, password string) (*entities.User, error)
	DeleteUser(username string) error
	CreateSession(username string) (entities.Session, error)
	// Find a user referenced by session token,
	// if the session has not expired.
	GetUserByToken(token string) (entities.User, entities.Session, error)
	DeleteSession(token string) error
}
