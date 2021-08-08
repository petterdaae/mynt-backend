package cmd

import (
	"fmt"
	internal "mynt/internal/routes"
	utils "mynt/internal/utils"
	"os"

	log "github.com/sirupsen/logrus"
)

// Execute is the main entry point of this application
func Execute() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{ForceColors: true})

	database := utils.InitDatabase()

	err := database.Ping()
	if err != nil {
		log.Fatal(fmt.Errorf("error occured while connecting to the database: %w", err))
	}

	err = database.Migrate()
	if err != nil {
		log.Fatal(fmt.Errorf("database migration failed: %w", err))
	}

	r := internal.SetupRoutes(database)

	r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}
