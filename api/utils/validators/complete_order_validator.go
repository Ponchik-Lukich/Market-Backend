package validators

import (
	"github.com/jmoiron/sqlx"
	"yandex-team.ru/bstask/api/models"
)

type ValidationCompleteOrderError struct {
	Message string
	Data    models.CompleteOrderDto
}

func (e *ValidationCompleteOrderError) Error() string {
	return e.Message
}

func ValidateCompleteOrder(db *sqlx.DB, order models.CompleteOrderDto) error {
	if order.OrderID <= 0 {
		return &ValidationCompleteOrderError{
			Message: "Order id is invalid",
			Data:    order,
		}
	}
	if order.CourierId <= 0 {
		return &ValidationCompleteOrderError{
			Message: "Courier id is invalid",
			Data:    order,
		}
	}
	if order.CompleteTime == nil {
		return &ValidationCompleteOrderError{
			Message: "Complete time is invalid",
			Data:    order,
		}
	}
	var assigned bool
	query := `SELECT assigned FROM orders WHERE id = $1`
	err := db.Get(&assigned, query, order.OrderID)
	if err != nil {
		return &ValidationCompleteOrderError{
			Message: "Order not found",
			Data:    order,
		}
	} else if assigned {
		return &ValidationCompleteOrderError{
			Message: "Order already completed",
			Data:    order,
		}
	}

	return nil
}
