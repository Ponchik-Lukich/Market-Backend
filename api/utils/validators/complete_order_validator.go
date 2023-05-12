package validators

import (
	"yandex-team.ru/bstask/api/models"
)

type ValidationCompleteOrderError struct {
	Message string
	Data    models.CompleteOrderDto
	Err     interface{}
}

func (e *ValidationCompleteOrderError) Error() string {
	return e.Message
}

func ValidateCompleteOrder(order models.CompleteOrderDto) error {
	if order.OrderID <= 0 {
		return &ValidationCompleteOrderError{
			Message: "Order id is invalid",
			Data:    order,
			Err:     order.OrderID,
		}
	}
	if order.CourierId <= 0 {
		return &ValidationCompleteOrderError{
			Message: "Courier id is invalid",
			Data:    order,
			Err:     order.CourierId,
		}
	}
	if order.CompleteTime == nil {
		return &ValidationCompleteOrderError{
			Message: "Complete time is invalid",
			Data:    order,
			Err:     order.CompleteTime,
		}
	}

	return nil
}

func ValidateIds(id int64, setIds *map[int64]struct{}) error {
	if _, ok := (*setIds)[id]; ok {
		return &ValidationCompleteOrderError{
			Message: "Completed order ids has duplicates",
		}
	}
	(*setIds)[id] = struct{}{}
	return nil
}

func ValidateAssignedOrders(err error, result int) error {
	if err != nil {
		return err
	}
	if result == 1 {
		return &ValidationCompleteOrderError{
			Message: "Orders courier_id is null",
		}
	}
	if result == 2 {
		return &ValidationCompleteOrderError{
			Message: "Orders is assigned to another courier",
		}
	}
	if result == 3 {
		return &ValidationCompleteOrderError{
			Message: "Order does not exist",
		}
	}
	return nil
}

func ValidateExistingCouriers(err error, ids []int64) error {
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil
		}
		return &ValidationCompleteOrderError{
			Message: "Data contains not existing couriers ids",
			Err:     ids,
		}
	}
	return nil
}
