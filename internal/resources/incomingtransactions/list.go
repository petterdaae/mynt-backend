package incomingtransactions

import (
	"backend/internal/types"
	"fmt"
)

func (resource *Resource) ListAll() ([]types.IncomingTransaction, error) {
	rows, err := resource.database.Query(`
		SELECT 
			id,
			account_id,
			accounting_date,
			interest_date,
			amount,
			text
    	FROM
        	incoming_transactions
		WHERE
			user_id = $1
    	ORDER BY accounting_date DESC
		`,
		resource.sub,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to query incoming_transactions: %w", err)
	}

	defer rows.Close()

	transactions := []types.IncomingTransaction{}
	for rows.Next() {
		var transaction types.IncomingTransaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.AccountID,
			&transaction.AccountingDate,
			&transaction.InterestDate,
			&transaction.Amount,
			&transaction.Text,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
