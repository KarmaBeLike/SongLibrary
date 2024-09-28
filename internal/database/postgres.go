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

// OpenDB открывает соединение с базой данных и выполняет ping
func OpenDB(cfg *config.Config) (*sql.DB, error) {
	dbURL := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Println("Failed to connect to the database:", err)
		return nil, err
	}

	// Пинг для проверки соединения
	if err := db.Ping(); err != nil {
		log.Println("Database ping failed:", err)
		return nil, err
	}

	log.Println("Database connection established.")
	return db, nil
}

// RunMigrations запускает миграции
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
	m.Up()

	log.Println("Migrations applied successfully!")
	return nil
}
