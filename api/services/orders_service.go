package services

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
	"yandex-team.ru/bstask/api/models"
	"yandex-team.ru/bstask/api/utils/validators"
)

func GetOrders(db *sqlx.DB, limit int, offset int) ([]models.Order, error) {
	var orders []models.Order
	query := `SELECT orders.id,
       orders.cost,
       orders.delivery_hours,
       orders.delivery_district,
       orders.weight,
       order_completion.completion_time
FROM orders
         FULL OUTER JOIN order_completion ON orders.id = order_completion.order_id
LIMIT $1 OFFSET $2`
	err := db.Select(&orders, query, limit, offset)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func GetOrderById(db *sqlx.DB, OrderID int64) (*models.Order, error) {
	var order models.Order
	query := `SELECT id, cost, delivery_hours, delivery_district, weight FROM orders WHERE id = $1`
	err := db.Get(&order, query, OrderID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

func CreateOrders(db *sqlx.DB, orders []models.CreateOrderDto) ([]models.Order, error) {
	var createdOrders []models.Order

	// Validate orders
	for _, order := range orders {
		err := validators.ValidateOrder(order)
		if err != nil {
			return nil, &validators.ValidationOrderError{
				Message: "Validation failed for order",
				Data:    order,
			}
		}
	}

	var query strings.Builder
	query.WriteString("INSERT INTO orders (cost, delivery_hours, delivery_district, weight) VALUES ")

	var values []interface{}
	for i, order := range orders {
		if i > 0 {
			query.WriteString(", ")
		}
		query.WriteString(fmt.Sprintf("($%d, $%d, $%d, $%d)", i*4+1, i*4+2, i*4+3, i*4+4))
		values = append(values, order.Cost, order.DeliveryHours, order.Regions, order.Weight)
	}

	query.WriteString(" RETURNING id, cost, delivery_hours, delivery_district, weight")
	rows, err := db.Query(query.String(), values...)
	if err != nil {
		panic(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var createdOrder models.Order
		err := rows.Scan(&createdOrder.OrderID, &createdOrder.Cost, &createdOrder.DeliveryHours, &createdOrder.Regions, &createdOrder.Weight)
		if err != nil {
			return nil, err
		}
		createdOrders = append(createdOrders, createdOrder)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return createdOrders, nil
}
