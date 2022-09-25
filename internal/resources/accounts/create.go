package accounts

import (
	"backend/internal/types"
	"fmt"
)

func (resource Resource) CreateIfNotExists(account *types.Account) error {
	err := resource.database.Exec(`
		INSERT INTO accounts (
			user_id,
			id,
			account_number,
			name,
			available,
			balance
		) VALUES (
			$1, $2, $3, $4, $5, $6
		) ON CONFLICT (id) DO UPDATE SET
			available = $5,
			balance = $6
		`,
		resource.sub,
		account.ID,
		account.AccountNumber,
		account.Name,
		account.Available,
		account.Balance,
	)

	if err != nil {
		return fmt.Errorf("failed to insert account: %w", err)
	}

	return nil
}
