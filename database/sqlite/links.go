package sqlite

import (
	"github.com/ordinary-dev/phoenix/database/entities"
)

func (db *SqliteDB) GetLinksFromGroup(groupID int) ([]entities.Link, error) {
	query := `
        SELECT id, name, href, group_id, icon
        FROM links
        WHERE group_id = ?
        ORDER BY id
    `

	rows, err := db.conn.Query(query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []entities.Link
	for rows.Next() {
		var link entities.Link
		if err := rows.Scan(&link.ID, &link.Name, &link.Href, &link.GroupID, &link.Icon); err != nil {
			return nil, err
		}
		links = append(links, link)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return links, nil
}

func (db *SqliteDB) GetLink(id int) (*entities.Link, error) {
	query := `
        SELECT id, name, href, group_id, icon
        FROM links
        WHERE id = ?
    `

	var link entities.Link
	err := db.conn.
		QueryRow(query, id).
		Scan(&link.ID, &link.Name, &link.Href, &link.GroupID, &link.Icon)
	if err != nil {
		return nil, err
	}

	return &link, nil
}

// Create a new link in the database.
// The function fills in the ID.
func (db *SqliteDB) CreateLink(link *entities.Link) error {
	query := `
        INSERT INTO links (name, href, group_id, icon)
        VALUES (?, ?, ?, ?)
        RETURNING id
    `

	err := db.conn.
		QueryRow(query, link.Name, link.Href, link.GroupID, link.Icon).
		Scan(&link.ID)

	if err != nil {
		return err
	}

	return nil
}

func (db *SqliteDB) UpdateLink(link *entities.Link) error {
	query := `
        UPDATE links
        SET name = ?, href = ?, group_id = ?, icon = ?
        WHERE id = ?
    `

	res, err := db.conn.Exec(query, link.Name, link.Href, link.GroupID, link.Icon, link.ID)
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

func (db *SqliteDB) DeleteLink(linkID int) error {
	query := `
        DELETE FROM links
        WHERE id = ?
    `

	res, err := db.conn.Exec(query, linkID)
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
