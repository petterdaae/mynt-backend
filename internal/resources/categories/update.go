package categories

import "backend/internal/types"

func (c *Resource) Update(id int64, draftCategory types.DraftCategory) error {
	err := c.database.Exec(`
		UPDATE categories 
		SET
			name = $3,
			parent_id = $4,
			color = $5,
			ignore = $6
		WHERE 
			user_id = $1 AND id = $2
		`,
		c.sub,
		id,
		draftCategory.Name,
		draftCategory.ParentID,
		draftCategory.Color,
		draftCategory.Ignore,
	)

	return err
}
