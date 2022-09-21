package user

import (
	"backend/internal/types"
	"fmt"
)

func (resource *Resource) Read() (types.User, error) {
	row, err := resource.database.QueryRow(`
		SELECT 
			sbanken_client_id,
			sbanken_client_secret
		FROM 
			users
		WHERE 
			sub = $1
		`,
		resource.sub,
	)

	if err != nil {
		return types.User{}, fmt.Errorf("failed to query users: %w", err)
	}

	var user types.User

	err = row.Scan(&user.SbankenClientID, &user.SbankenClientSecret)
	if err != nil {
		return types.User{}, fmt.Errorf("failed to query users: %w", err)
	}

	return user, nil
}
