package sql

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

type Transaction struct {
	TransactionId string    `json:"transaction_id"`
	Amount        float64   `json:"amount"`
	Spent         bool      `json:"spend"`
	CreatedAt     time.Time `json:"create_at"`
}

func OpenDatabase() (*sql.DB, error) {
	godotenv.Load()
	host := os.Getenv("DB_HOST")
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, fmt.Errorf("failed to convert DB_PORT to int: %v", err)
	}
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	postgresInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", postgresInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	return db, nil
}

func MigrateDatabase(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %v", err)
	}

	if err := m.Up(); err != migrate.ErrNoChange && err != nil {
		return fmt.Errorf("failed to apply migrations: %v", err)
	}

	return nil
}

func init() {
	var err error
	db, err = OpenDatabase()
	if err != nil {
		log.Fatalln(err)
	}

	defer db.Close()

	if err := MigrateDatabase(db); err != nil {
		log.Fatalln(err)
	}
}
