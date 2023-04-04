package validators

import (
	"github.com/jmoiron/sqlx"
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

func ValidateCompleteOrder(db *sqlx.DB, order models.CompleteOrderDto) error {
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
	var assigned bool
	query := `SELECT assigned FROM orders WHERE id = $1`
	err := db.Get(&assigned, query, order.OrderID)
	if err != nil {
		return &ValidationCompleteOrderError{
			Message: "Order not found",
			Data:    order,
			Err:     order.OrderID,
		}
	} else if assigned {
		return &ValidationCompleteOrderError{
			Message: "Order already completed",
			Data:    order,
			Err:     order.OrderID,
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
