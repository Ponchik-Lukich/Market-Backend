package validators

import "yandex-team.ru/bstask/api/models"

type ValidationOrderError struct {
	Message string
	Data    models.CreateOrderDto
}

func (e *ValidationOrderError) Error() string {
	return e.Message
}

// TODO add validation

func ValidateOrder(order models.CreateOrderDto) error {
	if order.Weight <= 0 {
		return &ValidationOrderError{
			Message: "Order weight is invalid",
			Data:    order,
		}
	}
	if order.Cost <= 0 {
		return &ValidationOrderError{
			Message: "Order cost is invalid",
			Data:    order,
		}
	}
	if len(order.DeliveryHours) == 0 {
		return &ValidationOrderError{
			Message: "Order delivery hours is empty",
			Data:    order,
		}
	}
	return nil
}
