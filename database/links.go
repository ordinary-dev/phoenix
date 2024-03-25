package database

type Link struct {
	ID      int
	Name    string
	Href    string
	GroupID int
	Icon    *string
}

func GetLinksFromGroup(groupID int) ([]Link, error) {
	query := `
        SELECT id, name, href, group_id, icon
        FROM links
        WHERE group_id = ?
        ORDER BY id
    `

	rows, err := DB.Query(query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []Link
	for rows.Next() {
		var link Link
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

func GetLink(id int) (*Link, error) {
	query := `
        SELECT id, name, href, group_id, icon
        FROM links
        WHERE id = ?
    `

	var link Link
	err := DB.
		QueryRow(query, id).
		Scan(&link.ID, &link.Name, &link.Href, &link.GroupID, &link.Icon)
	if err != nil {
		return nil, err
	}

	return &link, nil
}

// Create a new link in the database.
// The function fills in the ID.
func CreateLink(link *Link) error {
	query := `
        INSERT INTO links (name, href, group_id, icon)
        VALUES (?, ?, ?, ?)
        RETURNING id
    `

	err := DB.
		QueryRow(query, link.Name, link.Href, link.GroupID, link.Icon).
		Scan(&link.ID)

	if err != nil {
		return err
	}

	return nil
}

func UpdateLink(link *Link) error {
	query := `
        UPDATE links
        SET name = ?, href = ?, group_id = ?, icon = ?
        WHERE id = ?
    `

	res, err := DB.Exec(query, link.Name, link.Href, link.GroupID, link.Icon, link.ID)
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

func DeleteLink(linkID int) error {
	query := `
        DELETE FROM links
        WHERE id = ?
    `

	res, err := DB.Exec(query, linkID)
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
