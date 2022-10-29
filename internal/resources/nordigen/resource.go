package nordigen

import "backend/internal/utils"

type Resource struct {
	sub               string
	database          *utils.Database
	nordigenSecretID  string
	nordigenSecretKey string
}

func Configure(sub string, database *utils.Database, nordigenSecretID, nordigenSecretKey string) Resource {
	return Resource{
		sub:               sub,
		database:          database,
		nordigenSecretID:  nordigenSecretID,
		nordigenSecretKey: nordigenSecretKey,
	}
}
