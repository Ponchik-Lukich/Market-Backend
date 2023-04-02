package validators

import "yandex-team.ru/bstask/api/models"

type ValidationCompleteOrderError struct {
	Message string
	Data    models.CompleteOrderDto
}

func (e *ValidationCompleteOrderError) Error() string {
	return e.Message
}

// TODO add validation

func ValidateCompleteOrder(order models.CompleteOrderDto) error {
	if order.OrderID <= 0 {
		return &ValidationCompleteOrderError{
			Message: "Order weight is invalid",
			Data:    order,
		}
	}
	//if order.CompleteTime <= 0 {
	//	return &ValidationCompleteOrderError{
	//		Message: "Order cost is invalid",
	//		Data:    order,
	//	}
	//}
	//if len(order.DeliveryHours) == 0 {
	//	return &ValidationCompleteOrderError{
	//		Message: "Order delivery hours is empty",
	//		Data:    order,
	//	}
	//}
	return nil
}
