package utils

import (
	"database/sql"
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
		Box: packr.New("sql", "../../sql"),
	}
	n, err := migrate.Exec(c.DB, "postgres", migrations, migrate.Up)
	if err == nil {
		log.WithField("changes", n).Info("Successfully migrated database changes")
	}
	return err
}

func (d *Database) Query(query string, params ...interface{}) (*sql.Rows, error) {
	connection, err := d.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	defer connection.Close()

	rows, err := connection.Query(query, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to query database: %w", err)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("failed to query database: %w", rows.Err())
	}

	return rows, nil
}

func (d *Database) QueryRow(query string, params ...interface{}) (*sql.Row, error) {
	connection, err := d.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	defer connection.Close()

	return connection.QueryRow(query, params...), nil
}

func (d *Database) Exec(query string, params ...interface{}) error {
	connection, err := d.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer connection.Close()

	_, err = connection.Exec(query, params...)
	if err != nil {
		return fmt.Errorf("failed to query database: %w", err)
	}

	return nil
}
