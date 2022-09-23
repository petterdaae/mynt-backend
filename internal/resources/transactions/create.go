package transactions

import (
	"backend/internal/types"
	"fmt"
)

func (resource Resource) CreateIfNotExists(transaction *types.Transaction) error {
	err := resource.database.Exec(`
		INSERT INTO transactions (
			user_id,
			id,
			account_id,
			accounting_date,
			interest_date,
			amount,
			text
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		) ON CONFLICT DO NOTHING
		`,
		resource.sub,
		transaction.ID,
		transaction.AccountID,
		transaction.AccountingDate,
		transaction.InterestDate,
		transaction.Amount,
		transaction.Text,
	)

	if err != nil {
		return fmt.Errorf("failed to insert transaction: %w", err)
	}

	return nil
}
