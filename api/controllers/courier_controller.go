package controllers

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"yandex-team.ru/bstask/api/models"
	"yandex-team.ru/bstask/api/services"
	"yandex-team.ru/bstask/api/utils/validators"
)

func CreateCourier(c echo.Context, db *sqlx.DB) error {
	var req models.CreateCourierRequest
	var res models.CreateCourierResponse

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.BadRequestResponse{Error: "bad request"})
	}

	createdCouriers, err := services.CreateCouriers(db, req.Couriers)
	if err != nil {
		switch e := err.(type) {
		case *validators.ValidationCourierError:
			return c.JSON(http.StatusBadRequest, models.BadRequestResponse{
				Error:   fmt.Sprintf("Validation error for courier: %v", e.Data),
				Message: e.Message,
			})
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, models.InternalServerErrorResponse{
				Error: "Error creating couriers",
			})
		}
	}

	res.Couriers = createdCouriers
	return c.JSON(http.StatusOK, res)
}

func GetCourierById(c echo.Context, db *sqlx.DB) error {
	courierID, err := strconv.ParseInt(c.Param("courier_id"), 10, 64)
	if err != nil || courierID <= 0 {
		badRequest := models.BadRequestResponse{
			Error:   "bad request",
			Message: "courier_id must be a positive integer",
		}
		return c.JSON(http.StatusBadRequest, badRequest)
	}
	courier, err := services.GetCourierById(db, courierID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.InternalServerErrorResponse{
			Error: "Error getting courier",
		})
	}
	if courier == nil {
		return c.JSON(http.StatusNotFound, models.NotFoundResponse{Error: "not found"})
	}
	return c.JSON(http.StatusOK, courier)
}

func GetCouriers(c echo.Context, db *sqlx.DB) error {
	limitParam := c.QueryParam("limit")
	offsetParam := c.QueryParam("offset")
	var err error

	limit := 1
	offset := 0

	if limitParam != "" {
		limit, err = strconv.Atoi(limitParam)
		if err != nil || limit <= 0 {
			return c.JSON(http.StatusBadRequest, models.BadRequestResponse{
				Error:   "bad request",
				Message: "limit must be a positive integer"})
		}
	}
	if offsetParam != "" {
		offset, err = strconv.Atoi(offsetParam)
		if err != nil || offset < 0 {
			return c.JSON(http.StatusBadRequest, models.BadRequestResponse{
				Error:   "bad request",
				Message: "offset must be a positive integer"})
		}
	}

	couriers, err := services.GetCouriers(db, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.InternalServerErrorResponse{
			Error: "Error getting couriers",
		})
	}
	var res models.GetCourierResponse
	if len(couriers) == 0 {
		res.Couriers = []models.Courier{}
	} else {
		res.Couriers = couriers
	}
	res.Limit = limit
	res.Offset = offset

	return c.JSON(http.StatusOK, res)
}

//func CreateCourier(c echo.Context, db *sqlx.DB) error {
//	var req models.CreateCourierRequest
//	var res models.CreateCourierResponse
//	var couriers []models.CreateCourierDto
//
//	if err := c.Bind(&req); err != nil {
//		return echo.NewHTTPError(http.StatusBadRequest, models.BadRequestResponse{Error: "bad request"})
//	}
//	couriers = req.Couriers
//	err := services.CreateCouriers(db, couriers)
//	if err != nil {
//		switch e := err.(type) {
//		case *validators.ValidationCourierError:
//			return c.JSON(http.StatusBadRequest, models.BadRequestResponse{
//				Error: fmt.Sprintf("Validation error for courier: %v", e.Data),
//			})
//		default:
//			return echo.NewHTTPError(http.StatusInternalServerError, "Error creating couriers")
//		}
//	}
//	// return array of created couriers as json
//	return c.JSON(http.StatusOK, couriers)
//}
