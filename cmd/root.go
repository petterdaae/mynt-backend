package cmd

import (
	"fmt"
	"mynt/internal"
	"os"

	log "github.com/sirupsen/logrus"
)

// Execute is the main entry point of this application
func Execute() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{ForceColors: true})

	d := internal.InitDependencies()

	err := d.Db.Ping()
	if err != nil {
		log.Fatal(fmt.Errorf("error occured while connecting to the database: %w", err))
	}

	err = d.Db.Migrate()
	if err != nil {
		log.Fatal(fmt.Errorf("database migration failed: %w", err))
	}

	r := internal.SetupRoutes(d)

	r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}
