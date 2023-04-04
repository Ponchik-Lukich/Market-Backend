package services

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"strings"
	"yandex-team.ru/bstask/api/models"
	"yandex-team.ru/bstask/api/utils/validators"
)

func GetOrders(db *sqlx.DB, limit int, offset int) ([]models.Order, error) {
	var orders []models.Order
	query := `SELECT orders.id, orders.cost, orders.delivery_hours, orders.delivery_district, orders.weight, order_completion.complete_time
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
			return nil, err
		}
	}

	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	chunkSize := 16382
	for i := 0; i < len(orders); i += chunkSize {
		end := i + chunkSize
		if end > len(orders) {
			end = len(orders)
		}
		chunk := orders[i:end]

		var query strings.Builder
		query.WriteString("INSERT INTO orders (cost, delivery_hours, delivery_district, weight) VALUES ")

		var placeholders []string
		var values []interface{}
		for i, order := range chunk {
			placeholder := fmt.Sprintf("($%d, $%d, $%d, $%d)", i*4+1, i*4+2, i*4+3, i*4+4)
			placeholders = append(placeholders, placeholder)
			values = append(values, order.Cost, order.DeliveryHours, order.Regions, order.Weight)
		}

		query.WriteString(strings.Join(placeholders, ", "))
		query.WriteString(" RETURNING id, cost, delivery_hours, delivery_district, weight")

		rows, err := tx.Query(query.String(), values...)
		if err != nil {
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
	}

	return createdOrders, nil
}

func CompleteOrder(db *sqlx.DB, orders []models.CompleteOrderDto) ([]models.Order, error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	chunkSize := 21845

	var completedOrders []models.Order
	setIDs := map[int64]struct{}{}
	var orderIDs []int64
	for _, order := range orders {
		if err := validators.ValidateCompleteOrder(db, order); err != nil {
			return nil, err
		}
		if err := validators.ValidateIds(order.OrderID, &setIDs); err != nil {
			return nil, err
		}
		orderIDs = append(orderIDs, order.OrderID)
	}
	for i := 0; i < len(orderIDs); i += chunkSize {
		end := i + chunkSize
		if end > len(orders) {
			end = len(orders)
		}
		chunk := orderIDs[i:end]
		var placeholders []interface{}
		for _, id := range chunk {
			placeholders = append(placeholders, id)
		}
		query := "SELECT NOT EXISTS (" +
			"SELECT 1 FROM orders WHERE id IN ( " +
			strings.Trim(strings.Repeat("?, ", len(placeholders)), " ,") +
			") AND assignment = true)"
		var result bool
		println(query)
		err := tx.Get(&result, query, placeholders...)
		if err != nil {
			return nil, err
		}
		if result {
			return nil, errors.New("data contains completed orders")
		}
	}

	for i := 0; i < len(orders); i += chunkSize {
		end := i + chunkSize
		if end > len(orders) {
			end = len(orders)
		}
		chunk := orders[i:end]

		// Insert into order_completion table
		var placeholders []string
		var values []interface{}
		for i, order := range chunk {
			placeholder := fmt.Sprintf("($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3)
			placeholders = append(placeholders, placeholder)
			values = append(values, order.OrderID, order.CourierId, order.CompleteTime)
		}
		query := "INSERT INTO order_completion (order_id, courier_id, complete_time) VALUES " + strings.Join(placeholders, ", ")
		if _, err := tx.Exec(query, values...); err != nil {
			return nil, err
		}

		// Update assigned field in orders table
		var orderIDs []int64
		for _, order := range chunk {
			orderIDs = append(orderIDs, order.OrderID)
		}
		query = "UPDATE orders SET assigned = true WHERE id = ANY($1) RETURNING id, cost, delivery_hours, delivery_district, weight"
		rows, err := tx.Queryx(query, pq.Array(orderIDs))
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var completedOrder models.Order
			if err := rows.StructScan(&completedOrder); err != nil {
				return nil, err
			}
			for _, order := range chunk {
				if order.OrderID == completedOrder.OrderID {
					completedOrder.CompleteTime = order.CompleteTime
					break
				}
			}
			completedOrders = append(completedOrders, completedOrder)
		}
		if err := rows.Err(); err != nil {
			return nil, err
		}
	}

	return completedOrders, nil
}

//func CompleteOrder(db *sqlx.DB, orders []models.CompleteOrderDto) ([]models.Order, error) {
//	var completedOrders []models.Order
//	for _, order := range orders {
//		if err := validators.ValidateCompleteOrder(order); err != nil {
//			return nil, &validators.ValidationCompleteOrderError{
//				Message: "Validation failed for completed order",
//				Data:    order,
//			}
//		}
//	}
//
//	tx, err := db.Beginx()
//	if err != nil {
//		return nil, err
//	}
//	defer func() {
//		if p := recover(); p != nil {
//			tx.Rollback()
//			panic(p)
//		} else if err != nil {
//			tx.Rollback()
//		} else {
//			tx.Commit()
//		}
//	}()
//
//	// Insert into order_completion table
//	query := `INSERT INTO order_completion (order_id, courier_id, complete_time) VALUES ($1, $2, $3)`
//	for _, order := range orders {
//		if _, err := tx.Exec(query, order.OrderID, order.CourierId, order.CompleteTime); err != nil {
//			return nil, err
//		}
//	}
//
//	// Update assigned field in orders table
//	query = `UPDATE orders SET assigned = true WHERE id = ANY($1) RETURNING id, cost, delivery_hours, delivery_district, weight`
//	var orderIDs []int64
//	for _, order := range orders {
//		orderIDs = append(orderIDs, order.OrderID)
//	}
//	rows, err := tx.Queryx(query, pq.Array(orderIDs))
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//
//	for rows.Next() {
//		var completedOrder models.Order
//		if err := rows.StructScan(&completedOrder); err != nil {
//			return nil, err
//		}
//		for _, order := range orders {
//			if order.OrderID == completedOrder.OrderID {
//				completedOrder.CompleteTime = order.CompleteTime
//				break
//			}
//		}
//		completedOrders = append(completedOrders, completedOrder)
//	}
//
//	return completedOrders, nil
//}
