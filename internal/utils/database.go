package utils

import (
	"fmt"
	"os"

	"github.com/gobuffalo/packr/v2"
	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"
)

type Database struct {
	psqlConnectionString string
}

func InitDatabase() *Database {
	return &Database{
		psqlConnectionString: fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			os.Getenv("POSTGRES_HOST"),
			os.Getenv("POSTGRES_PORT"),
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"),
			os.Getenv("POSTGRES_DB"),
			os.Getenv("POSTGRES_SSL"),
		),
	}
}

func (d *Database) Connect() (*sqlx.DB, error) {
	return sqlx.Connect("postgres", d.psqlConnectionString)
}

func (d *Database) Ping() error {
	c, err := d.Connect()
	if err != nil {
		return err
	}
	err = c.Ping()
	return err
}

func (d *Database) Migrate() error {
	c, err := d.Connect()
	if err != nil {
		return err
	}
	migrations := &migrate.PackrMigrationSource{
		Box: packr.New("sql", "../sql"),
	}
	n, err := migrate.Exec(c.DB, "postgres", migrations, migrate.Up)
	if err == nil {
		log.WithField("changes", n).Info("Successfully migrated database changes")
	}
	return err
}
