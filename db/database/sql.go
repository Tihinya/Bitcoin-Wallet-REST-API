package sql

import (
	"bitcoin-wallet/config"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/caarlos0/env"
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
	postgresInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.EnvConfig.Host, config.EnvConfig.Port, config.EnvConfig.User, config.EnvConfig.Password, config.EnvConfig.DbName)

	db, err := sql.Open("postgres", postgresInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	return db, nil
}

func MigrateDatabase(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{DatabaseName: config.EnvConfig.DbName, MigrationsTable: postgres.DefaultMigrationsTable, MultiStatementEnabled: false, MultiStatementMaxSize: postgres.DefaultMultiStatementMaxSize})
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
	err = godotenv.Load(".env")
	if err != nil {
		log.Println(err)
	}
	err = env.Parse(&config.EnvConfig)
	if err != nil {
		log.Println(err)
	}
	log.Println(config.EnvConfig)
	db, err = OpenDatabase()
	if err != nil {
		log.Fatalln(err)
	}

	// if err := MigrateDatabase(db); err != nil {
	// 	log.Fatalln(err)
	// }
}
