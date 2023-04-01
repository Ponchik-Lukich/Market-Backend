package controllers

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"yandex-team.ru/bstask/api/models"
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

	couriers, err := services.GetCouriers(db, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error getting couriers")
	}
	var response models.GetCouriersResponse
	response.Couriers = couriers
	response.Limit = limit
	response.Offset = offset

	return c.JSON(http.StatusOK, response)
}
