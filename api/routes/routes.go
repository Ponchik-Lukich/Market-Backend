package routes

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"yandex-team.ru/bstask/api/controllers"
)

func SetupRoutes(e *echo.Echo, db *sqlx.DB) {
	// Couriers
	e.POST("/couriers", controllers.CreateCourier)
	e.GET("/couriers/:courier_id", controllers.GetCourier)
	e.GET("/couriers", func(c echo.Context) error {
		return controllers.GetCouriers(c, db)
	})
	// Orders
	e.POST("/orders", controllers.CreateOrder)
	e.GET("/orders/:order_id", controllers.GetOrder)
	e.GET("/orders", controllers.GetOrders)
	e.POST("/orders/complete", controllers.CompleteOrder)

	// Ping
	e.GET("/ping", ping)
}
