package transactions

import "backend/internal/utils"

type Resource struct {
	sub      string
	database *utils.Database
}

func Configure(sub string, database *utils.Database) Resource {
	return Resource{
		sub:      sub,
		database: database,
	}
}
