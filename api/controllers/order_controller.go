package controllers

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"yandex-team.ru/bstask/api/models"
	"yandex-team.ru/bstask/api/services"
)

//TODO: Implement logic for order controller

func CreateOrder(c echo.Context) error {
	// Implement logic for creating an order
	return c.String(200, "Order created")
}

func GetOrder(c echo.Context, db *sqlx.DB) error {
	orderID, err := strconv.ParseInt(c.Param("order_id"), 10, 64)
	if err != nil || orderID <= 0 {
		badRequest := models.BadRequestResponse{Error: "bad request"}
		return c.JSON(http.StatusBadRequest, badRequest)
	}
	order, err := services.GetOrderById(db, orderID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.InternalServerErrorResponse{
			Error: "Error getting order",
		})
	}
	if order == nil {
		return c.JSON(http.StatusNotFound, models.NotFoundResponse{Error: "not found"})
	}
	return c.JSON(http.StatusOK, order)
}

func GetOrders(c echo.Context, db *sqlx.DB) error {
	limitParam := c.QueryParam("limit")
	offsetParam := c.QueryParam("offset")

	limit := 10
	offset := 0

	if limitParam != "" {
		limit, err := strconv.Atoi(limitParam)
		if err != nil || limit <= 0 {
			badRequest := models.BadRequestResponse{Error: "bad request"}
			return c.JSON(http.StatusBadRequest, badRequest)
		}
	}
	if offsetParam != "" {
		offset, err := strconv.Atoi(offsetParam)
		if err != nil || offset < 0 {
			badRequest := models.BadRequestResponse{Error: "bad request"}
			return c.JSON(http.StatusBadRequest, badRequest)
		}
	}

	orders, err := services.GetOrders(db, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.InternalServerErrorResponse{
			Error: "Error getting orders",
		})
	}
	var res models.GetOrderResponse
	res.Orders = orders
	res.Limit = limit
	res.Offset = offset

	return c.JSON(http.StatusOK, res)
}

func CompleteOrder(c echo.Context) error {
	// Implement logic for completing an order
	return c.String(200, "Order completed")
}
