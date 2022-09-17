package categories

func (c *Resource) Delete(id int64) error {
	err := c.database.Exec(`
		DELETE FROM categories 
		WHERE 
			user_id = $1 AND id = $2
		`,
		c.sub,
		id,
	)

	return err
}
