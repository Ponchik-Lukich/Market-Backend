package validators

import "yandex-team.ru/bstask/api/models"

type ValidationOrderError struct {
	Message string
	Data    models.CreateOrderDto
}

func (e *ValidationOrderError) Error() string {
	return e.Message
}

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
	set := map[string]struct{}{}
	for i := 0; i < len(order.DeliveryHours); i++ {
		if _, ok := set[order.DeliveryHours[i]]; ok {
			return &ValidationOrderError{
				Message: "Order delivery hours has duplicates",
				Data:    order,
			}
		} else {
			set[order.DeliveryHours[i]] = struct{}{}
		}
		if !ValidateTime(order.DeliveryHours[i]) {
			return &ValidationOrderError{
				Message: "Order delivery hour is invalid",
				Data:    order,
			}
		}
	}
	return nil
}
