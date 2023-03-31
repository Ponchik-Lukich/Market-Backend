package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
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
	err := config.InitializeTables(db)
	if err != nil {
		panic(err)
	}

	e := setupServer()
	e.Logger.Fatal(e.Start(":8080"))
}

func setupServer() *echo.Echo {
	e := echo.New()
	routes.SetupRoutes(e)
	return e
}
