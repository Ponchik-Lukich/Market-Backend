package config

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "postgres"
)

func ConnectDB() *sqlx.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}

	log.Println("Successfully connected to the database")
	return db

}

func InitializeTables(db *sqlx.DB) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	folder := fmt.Sprintf("%s/%s/%s/%s", dir, "api", "models", "queries")
	tables := []string{"couriers", "orders", "courier_work_shifts", "order_completion"}

	for _, table := range tables {
		sql, err := os.ReadFile(fmt.Sprintf("%s/%s.sql", folder, table))

		if err != nil {
			return err
		}

		_, err = db.Exec(string(sql))

		if err != nil {
			return err
		}
	}
	return nil
}

func DropTables(db *sqlx.DB) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	folder := fmt.Sprintf("%s/%s/%s/%s", dir, "api", "models", "queries")
	sql, err := os.ReadFile(fmt.Sprintf("%s/%s.sql", folder, "down"))
	if err != nil {
		return err
	}
	_, err = db.Exec(string(sql))

	if err != nil {
		return err
	}
	return nil
}
