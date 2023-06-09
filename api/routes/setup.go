package routes

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"yandex-team.ru/bstask/api/controllers"
)

func SetupRoutes(e *echo.Echo, db *sqlx.DB, rateLimit float64, rateBurst int) {
	// Couriers
	routesLimiter := newLimiter(rateLimit, rateBurst)
	e.POST("/couriers", withLimiter(routesLimiter, controllers.CreateCourier, db))
	e.GET("/couriers/:courier_id", withLimiter(routesLimiter, controllers.GetCourierById, db))
	e.GET("/couriers", withLimiter(routesLimiter, controllers.GetCouriers, db))
	e.GET("/couriers/meta-info/:courier_id", withLimiter(routesLimiter, controllers.GetCourierMetaInfo, db))

	// Orders
	e.POST("/orders", withLimiter(routesLimiter, controllers.CreateOrder, db))
	e.GET("/orders/:order_id", withLimiter(routesLimiter, controllers.GetOrder, db))
	e.GET("/orders", withLimiter(routesLimiter, controllers.GetOrders, db))
	e.POST("/orders/complete", withLimiter(routesLimiter, controllers.CompleteOrder, db))

	e.GET("/ping", ping)
	// For testing because i didn't implement assignment
	e.GET("/test/:courier_id", withLimiter(routesLimiter, controllers.UpdateCourier, db))
}
