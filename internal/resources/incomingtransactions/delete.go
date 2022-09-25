package incomingtransactions

func (resource *Resource) DeleteAll(accountID string) error {
	return resource.database.Exec(`
		DELETE FROM incoming_transactions 
		WHERE user_id = $1 AND account_id = $2
		`,
		resource.sub,
		accountID,
	)
}
