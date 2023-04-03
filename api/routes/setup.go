package routes

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"yandex-team.ru/bstask/api/controllers"
)

func SetupRoutes(e *echo.Echo, db *sqlx.DB) {
	// Couriers
	routesLimiter := newLimiter(10, 1)
	e.POST("/couriers", withLimiter(routesLimiter, controllers.CreateCourier, db))
	e.GET("/couriers/:courier_id", withLimiter(routesLimiter, controllers.GetCourierById, db))
	e.GET("/couriers", withLimiter(routesLimiter, controllers.GetCouriers, db))
	e.GET("/couriers/meta-info/:courier_id", withLimiter(routesLimiter, controllers.GetCourierMetaInfo, db))

	// Orders
	e.POST("/orders", withLimiter(routesLimiter, controllers.CreateOrder, db))
	e.GET("/orders/:order_id", withLimiter(routesLimiter, controllers.GetOrder, db))
	e.GET("/orders", withLimiter(routesLimiter, controllers.GetOrders, db))
	e.POST("/orders/complete", withLimiter(routesLimiter, controllers.CompleteOrder, db))

	// Ping
	e.GET("/ping", ping)
}
