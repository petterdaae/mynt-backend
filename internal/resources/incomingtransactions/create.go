package incomingtransactions

import (
	"backend/internal/types"
	"fmt"
)

func (resource *Resource) Create(transaction *types.DraftIncomingTransaction) (int64, error) {
	row, err := resource.database.QueryRow(`
		INSERT INTO incoming_transactions (
			user_id,
			account_id,
			accounting_date,
			interest_date,
			amount,
			text
		) VALUES (
			$1, $2, $3, $4, $5, $6
		) RETURNING id
		`,
		resource.sub,
		transaction.AccountID,
		transaction.AccountingDate,
		transaction.InterestDate,
		transaction.Amount,
		transaction.Text,
	)

	if err != nil {
		return -1, fmt.Errorf("failed to insert incoming transaction: %w", err)
	}

	var id int64
	err = row.Scan(&id)

	if err != nil {
		return -1, fmt.Errorf("failed to insert category: %w", err)
	}

	return id, nil
}
