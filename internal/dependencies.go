package internal

// Dependencies contain common resources that the different routes need
type Dependencies struct {
	Db *database
}

// InitDependencies reads the environment and initializes the associated resources
func InitDependencies() *Dependencies {
	return &Dependencies{
		Db: initDatabase(),
	}
}
