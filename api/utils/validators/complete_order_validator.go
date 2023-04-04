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
			Message: "Data contains not existing orders ids",
		}
	}
	if result == 2 {
		return &ValidationCompleteOrderError{
			Message: "Data contains completed orders",
		}
	}
	return nil
}

func ValidateExistingCouriers(err error, result bool) error {
	if err != nil {
		return err
	}
	if !result {
		return &ValidationCompleteOrderError{
			Message: "Data contains not existing couriers ids",
		}
	}
	return nil
}
