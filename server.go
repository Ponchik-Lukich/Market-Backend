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
	//areas := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	//hours := []string{"10:00-12:00", "13:00-15:00", "16:00-18:00", "19:00-21:00"}
	////insert in couriers table
	//query := `INSERT INTO couriers (id, type, working_areas, working_hours) VALUES ($1, $2, $3, $4)`
	//_, err = db.Exec(query, 1, models.FOOT, pq.Array(areas), pq.Array(hours))
	//if err != nil {
	//	panic(err)
	//}
	e := setupServer()
	e.Logger.Fatal(e.Start(":8080"))
}

func setupServer() *echo.Echo {
	e := echo.New()
	routes.SetupRoutes(e, db)
	return e
}
