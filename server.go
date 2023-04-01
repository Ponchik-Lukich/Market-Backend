package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"yandex-team.ru/bstask/api/routes"
	"yandex-team.ru/bstask/config"
)

var db *sqlx.DB

func main() {
	db = config.ConnectDB()
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)
	err := config.DropTables(db)
	if err != nil {
		panic(err)
	}
	err = config.InitializeTables(db)
	if err != nil {
		panic(err)
	}
	hours := []string{"10:00-12:00", "13:00-15:00", "16:00-18:00", "19:00-21:00"}
	query := `INSERT INTO orders (id, cost, weight, delivery_hours, delivery_district) VALUES ($1, $2, $3, $4, $5)`
	_, err = db.Exec(query, 1, 100, 5, pq.Array(hours), 333)
	if err != nil {
		panic(err)
	}
	e := setupServer()
	e.Logger.Fatal(e.Start(":8080"))
}

func setupServer() *echo.Echo {
	e := echo.New()
	routes.SetupRoutes(e, db)
	return e
}
