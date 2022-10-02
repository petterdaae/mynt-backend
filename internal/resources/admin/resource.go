package admin

import "backend/internal/utils"

type Resource struct {
	database *utils.Database
}

func Configure(database *utils.Database) Resource {
	return Resource{
		database: database,
	}
}
