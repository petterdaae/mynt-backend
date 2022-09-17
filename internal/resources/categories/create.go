package categories

import (
	"backend/internal/types"
	"fmt"
)

func (c *Resource) Create(draftCategory types.DraftCategory) (int64, error) {
	row, err := c.database.QueryRow(`
		INSERT INTO categories (
			user_id,
			name,
			parent_id,
			color,
			ignore
		) VALUES (
			$1, $2, $3, $4, $5
		) RETURNING id
		`,
		c.sub,
		draftCategory.Name,
		draftCategory.ParentID,
		draftCategory.Color,
		draftCategory.Ignore,
	)

	if err != nil {
		return -1, fmt.Errorf("failed to insert category: %w", err)
	}

	var id int64
	err = row.Scan(&id)

	if err != nil {
		return -1, fmt.Errorf("failed to insert category: %w", err)
	}

	return id, err
}
