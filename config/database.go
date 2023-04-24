package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func ConnectDB() *sqlx.DB {
	dbHost := os.Getenv("DB_HOST")
	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatalf("Failed to parse the database port: %v", err)
	}
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)

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
	tables := []string{"couriers", "orders"}

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
