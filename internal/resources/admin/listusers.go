package admin

import (
	"fmt"
)

func (resource *Resource) ListUserIds() ([]string, error) {
	rows, err := resource.database.Query(`
		SELECT 
			id
		FROM 
			users
		`,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}

	defer rows.Close()

	var usersIDs []string
	for rows.Next() {
		var userID string
		err = rows.Scan(&userID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		usersIDs = append(usersIDs, userID)
	}

	return usersIDs, nil
}
