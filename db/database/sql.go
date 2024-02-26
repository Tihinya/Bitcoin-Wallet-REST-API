package sql

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "secret"
	dbname   = "postgres"
)

var db *sql.DB

type Transaction struct {
	TransactionId string    `json:"transaction_id"`
	Amount        float64   `json:"amount"`
	Spent         bool      `json:"spend"`
	CreatedAt     time.Time `json:"create_at"`
}

func openDatabase() *sql.DB {
	fmt.Println("open db")

	postgresInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", postgresInfo)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	return db
}

func makeMigration() {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Println(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres", driver)

	if err != nil {
		log.Println(err)
	}

	err = m.Up()
	if err != migrate.ErrNoChange && err != nil {
		log.Fatalln(err)
	}
}

func init() {
	db = openDatabase()
	makeMigration()
}
