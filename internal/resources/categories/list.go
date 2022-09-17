package categories

import (
	"backend/internal/types"
	"fmt"
)

func (c *Resource) List() ([]types.Category, error) {
	rows, err := c.database.Query(`
		SELECT 
			id,
			name,
			parent_id,
			color,
			ignore
		FROM 
			categories
		WHERE 
			user_id = $1
		ORDER BY 
			name
		`,
		c.sub,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to query categories: %w", err)
	}

	defer rows.Close()

	categories := []types.Category{}
	for rows.Next() {
		var category types.Category
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.ParentID,
			&category.Color,
			&category.Ignore,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		categories = append(categories, category)
	}

	return categories, nil
}
