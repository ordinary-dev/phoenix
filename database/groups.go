package database

type Group struct {
	ID    int
	Name  string
	Links []Link
}

func GetGroupsWithLinks() ([]Group, error) {
	query := `
        SELECT id, name
        FROM groups
        ORDER BY groups.id
    `

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []Group
	for rows.Next() {
		var group Group
		if err := rows.Scan(&group.ID, &group.Name); err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i := range groups {
		groups[i].Links, err = GetLinksFromGroup(groups[i].ID)
		if err != nil {
			return nil, err
		}
	}

	return groups, nil
}

// Create a new group in the database.
// The function fills in the ID.
func CreateGroup(group *Group) error {
	query := `
        INSERT INTO groups (name)
        VALUES (?)
        RETURNING id
    `

	if err := DB.QueryRow(query, group.Name).Scan(&group.ID); err != nil {
		return err
	}

	return nil
}

func UpdateGroup(id int, name string) error {
	query := `
        UPDATE groups
        SET name = ?
        WHERE id = ?
    `

	res, err := DB.Exec(query, name, id)
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

func DeleteGroup(groupID int) error {
	query := `
        DELETE FROM groups
        WHERE id = ?
    `

	res, err := DB.Exec(query, groupID)
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
