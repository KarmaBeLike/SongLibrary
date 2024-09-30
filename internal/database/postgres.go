package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/KarmaBeLike/SongLibrary/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

func OpenDB(cfg *config.Config) (*sql.DB, error) {
	dbURL := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to connect to the database:")
	}

	if err := db.Ping(); err != nil {
		return nil, errors.Wrap(err, "Database ping failed")
	}

	log.Println("Database connection established.")
	return db, nil
}

func RunMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return errors.Wrap(err, "create driver")
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		return errors.Wrap(err, "migrate")
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return errors.Wrap(err, "failed to apply migrations")
	}

	log.Println("Migrations applied successfully!")

	if err := printExistingTables(db); err != nil {
		return errors.Wrap(err, "failed to print existing tables")
	}

	return nil
}

func printExistingTables(db *sql.DB) error {
	rows, err := db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema='public'")
	if err != nil {
		return errors.Wrap(err, "querying existing tables")
	}
	defer rows.Close()

	log.Println("Existing tables:")
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return errors.Wrap(err, "scanning table name")
		}
		log.Println("-", tableName)
	}

	if err := rows.Err(); err != nil {
		return errors.Wrap(err, "error during rows iteration")
	}

	return nil
}
