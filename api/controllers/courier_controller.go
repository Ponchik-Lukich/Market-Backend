package controllers

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"yandex-team.ru/bstask/api/services"
)

//TODO: Implement logic for courier controller

func CreateCourier(c echo.Context) error {
	// Implement logic for creating a courier

	return c.String(200, "Courier created")
}

func GetCourier(c echo.Context) error {
	// Implement logic for getting a specific courier
	return c.String(200, "Courier found")
}

//func GetCouriers(c echo.Context) error {
//	// Implement logic for getting all couriers
//	// return couriers as json
//	couriers, err := services.GetCouriers(db, 10, 0)
//}

func GetCouriers(c echo.Context, db *sqlx.DB) error {
	//db := c.Get("db").(*sqlx.DB)

	// Get the limit and offset parameters if they are provided
	limitParam := c.QueryParam("limit")
	offsetParam := c.QueryParam("offset")

	// Default values if limit and offset are not provided
	limit := 10
	offset := 0

	if limitParam != "" {
		limit, _ = strconv.Atoi(limitParam)
	}
	if offsetParam != "" {
		offset, _ = strconv.Atoi(offsetParam)
	}

	// Call the service function to get couriers
	couriers, err := services.GetCouriers(db, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error getting couriers")
	}

	// Return the couriers as JSON
	return c.JSON(http.StatusOK, couriers)
}
