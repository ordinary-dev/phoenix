package sqlite

import (
	"github.com/ordinary-dev/phoenix/database/entities"
)

func (db *SqliteDB) GetGroupsWithLinks(username *string) ([]entities.Group, error) {
	query := `
        SELECT id, name
        FROM groups
        WHERE username = ?
        ORDER BY groups.id
    `

	rows, err := db.conn.Query(query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []entities.Group
	for rows.Next() {
		group := entities.Group{
			Username: username,
		}
		if err := rows.Scan(&group.ID, &group.Name); err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i := range groups {
		groups[i].Links, err = db.GetLinksFromGroup(groups[i].ID)
		if err != nil {
			return nil, err
		}
	}

	return groups, nil
}

// Create a new group in the database.
// The function fills in the ID.
func (db *SqliteDB) CreateGroup(group *entities.Group) error {
	query := `
        INSERT INTO groups (name, username)
        VALUES (?, ?)
        RETURNING id
    `

	if err := db.conn.QueryRow(query, group.Name, group.Username).Scan(&group.ID); err != nil {
		return err
	}

	return nil
}

// Get group by id without child links.
func (db *SqliteDB) GetGroup(id int) (entities.Group, error) {
	query := `
        SELECT name, username
        FROM groups
        WHERE id = ?
    `

	group := entities.Group{
		ID: id,
	}
	err := db.conn.QueryRow(query, id).Scan(&group.Name, &group.Username)
	return group, err
}

// Transfer groups from one owner to another.
func (db *SqliteDB) TransferGroups(from, to *string) error {
	query := `
        UPDATE groups
        SET username = ?
        WHERE username = ?
    `
	_, err := db.conn.Exec(query, to, from)
	return err
}

func (db *SqliteDB) UpdateGroup(id int, name string) error {
	query := `
        UPDATE groups
        SET name = ?
        WHERE id = ?
    `

	res, err := db.conn.Exec(query, name, id)
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

func (db *SqliteDB) DeleteGroup(groupID int) error {
	query := `
        DELETE FROM groups
        WHERE id = ?
    `

	res, err := db.conn.Exec(query, groupID)
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
