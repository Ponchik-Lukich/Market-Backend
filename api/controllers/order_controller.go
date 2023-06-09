package controllers

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"net/http"
	"reflect"
	"strconv"
	"yandex-team.ru/bstask/api/models"
	"yandex-team.ru/bstask/api/services"
	"yandex-team.ru/bstask/api/utils/validators"
)

func CreateOrder(c echo.Context, db *sqlx.DB) error {
	var req models.CreateOrderRequest
	var res models.CreateOrderResponse

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.BadRequestResponse{Error: "bad request"})
	}

	createdOrders, err := services.CreateOrders(db, req.Orders)
	if err != nil {
		switch e := err.(type) {
		case *validators.ValidationOrderError:
			return c.JSON(http.StatusBadRequest, models.BadRequestResponse{
				Error:   fmt.Sprintf("Validation error for order"),
				Message: e.Message,
				Data:    fmt.Sprintf("%v -> %v", e.Data, e.Err),
			})
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, models.InternalServerErrorResponse{
				Error: "Error creating orders",
			})
		}
	}
	res.Orders = createdOrders
	return c.JSON(http.StatusOK, res)
}

func GetOrder(c echo.Context, db *sqlx.DB) error {
	orderID, err := strconv.ParseInt(c.Param("order_id"), 10, 64)
	if err != nil || orderID <= 0 {
		badRequest := models.BadRequestResponse{
			Error:   "Bad request",
			Message: "order_id must be a positive integer",
		}
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
	var err error
	limit := 1
	offset := 0

	if limitParam != "" {
		limit, err = strconv.Atoi(limitParam)
		if err != nil || limit <= 0 {
			badRequest := models.BadRequestResponse{
				Error:   "bad request",
				Message: "limit must be a positive integer"}
			return c.JSON(http.StatusBadRequest, badRequest)
		}
	}
	if offsetParam != "" {
		offset, err = strconv.Atoi(offsetParam)
		if err != nil || offset < 0 {
			badRequest := models.BadRequestResponse{
				Error:   "bad request",
				Message: "offset must be a positive integer"}
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
	if len(orders) == 0 {
		res.Orders = []models.Order{}
	} else {
		res.Orders = orders
	}
	res.Limit = limit
	res.Offset = offset

	return c.JSON(http.StatusOK, res)
}

func CompleteOrder(c echo.Context, db *sqlx.DB) error {
	var req models.CompleteOrderRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.BadRequestResponse{Error: "bad request"})
	}
	completedOrders, err := services.CompleteOrder(db, req.Orders)
	if err != nil {
		switch e := err.(type) {
		case *validators.ValidationCompleteOrderError:
			if e.Err == nil {
				return c.JSON(http.StatusBadRequest, models.BadRequestResponse{
					Error:   fmt.Sprintf("Validation error for order"),
					Message: e.Message,
				})
			} else if reflect.TypeOf(e.Err).Kind() == reflect.Slice {
				return c.JSON(http.StatusBadRequest, models.BadRequestResponse{
					Error:   fmt.Sprintf("Validation error for order"),
					Message: e.Message,
					Data:    fmt.Sprintf("%v", e.Err),
				})
			} else {
				return c.JSON(http.StatusBadRequest, models.BadRequestResponse{
					Error:   fmt.Sprintf("Validation error for order"),
					Message: e.Message,
					Data:    fmt.Sprintf("%v -> %v", e.Data, e.Err),
				})
			}

		default:
			return echo.NewHTTPError(http.StatusInternalServerError, models.InternalServerErrorResponse{
				Error: "Error creating completed orders",
			})
		}
	}
	return c.JSON(http.StatusOK, completedOrders)
}
