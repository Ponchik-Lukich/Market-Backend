package validators

import "yandex-team.ru/bstask/api/models"

type ValidationOrderError struct {
	Message string
	Data    models.CreateOrderDto
}

func (e *ValidationOrderError) Error() string {
	return e.Message
}

func ValidateOrder(courier models.CreateOrderDto) error {
	if courier.Weight <= 0 {
		return &ValidationOrderError{
			Message: "Order weight is invalid",
			Data:    courier,
		}
	}
	if courier.Cost <= 0 {
		return &ValidationOrderError{
			Message: "Order cost is invalid",
			Data:    courier,
		}
	}
	if len(courier.DeliveryHours) == 0 {
		return &ValidationOrderError{
			Message: "Order delivery hours is empty",
			Data:    courier,
		}
	}
	return nil
}
