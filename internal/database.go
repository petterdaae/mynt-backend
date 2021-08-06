package internal

import (
	"fmt"
	"os"

	"github.com/gobuffalo/packr/v2"
	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"
)

type database struct {
	psqlConnectionString string
}

func initDatabase() *database {
	return &database{
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

func (d *database) connect() (*sqlx.DB, error) {
	return sqlx.Connect("postgres", d.psqlConnectionString)
}

func (d *database) Ping() error {
	c, err := d.connect()
	if err != nil {
		return err
	}
	err = c.Ping()
	return err
}

func (d *database) Migrate() error {
	c, err := d.connect()
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
