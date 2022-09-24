package transactions

import (
	"backend/internal/types"
	"fmt"
)

func (resource *Resource) List(from, to string) ([]types.Transaction, error) {
	rows, err := resource.database.Query(`
		SELECT 
        	t.id as id, 
        	t.account_id as account_id, 
        	split_part(t.accounting_date, 'T', 1) as accounting_date,
        	split_part(t.interest_date, 'T', 1) as interest_date,
        	t.custom_date as custom_date,
        	t.amount as amount,
        	t.text as text
    	FROM
        	transactions as t
    	WHERE
				t.user_id = $1
        	AND split_part(coalesce(t.custom_date, t.accounting_date), 'T', 1) >= $2
        	AND split_part(coalesce(t.custom_date, t.accounting_date), 'T', 1) <= $3
    	ORDER BY
        	coalesce(t.custom_date, t.accounting_date) desc, t.id;
		`,
		resource.sub,
		from,
		to,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}

	defer rows.Close()

	transactions := []types.Transaction{}
	for rows.Next() {
		var transaction types.Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.AccountID,
			&transaction.AccountingDate,
			&transaction.InterestDate,
			&transaction.CustomDate,
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
