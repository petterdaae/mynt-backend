package cmd

import (
	internal "backend/internal/routes"
	utils "backend/internal/utils"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

// Execute is the main entry point of this application
func Execute() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.JSONFormatter{})

	database := utils.InitDatabase()

	err := database.Ping()
	if err != nil {
		log.Fatal(fmt.Errorf("error occurred while connecting to the database: %w", err))
	}

	err = database.Migrate()
	if err != nil {
		log.Fatal(fmt.Errorf("database migration failed: %w", err))
	}

	r := internal.SetupRoutes(database)

	err = r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))

	if err != nil {
		panic(err)
	}
}
