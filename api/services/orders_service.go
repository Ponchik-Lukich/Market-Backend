package services

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"os"
	"strings"
	"yandex-team.ru/bstask/api/middleware"
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
	defer middleware.RollbackOrCommit(tx, &err)

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
	defer middleware.RollbackOrCommit(tx, &err)
	var completedOrders []models.Order
	setIDs := map[int64]struct{}{}
	var orderIDs, courierIDs []int64
	for _, order := range orders {
		if err := validators.ValidateCompleteOrder(order); err != nil {
			return nil, err
		}
		if err := validators.ValidateIds(order.OrderID, &setIDs); err != nil {
			return nil, err
		}
		orderIDs = append(orderIDs, order.OrderID)
		courierIDs = append(courierIDs, order.CourierId)
	}
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	folder := fmt.Sprintf("%s/%s/%s/%s", dir, "api", "models", "queries")
	files := [3]string{"check_orders", "check_couriers", "update_orders"}
	query, _ := os.ReadFile(fmt.Sprintf("%s/%s.sql", folder, files[0]))
	var result int
	var result1 bool
	err = tx.Get(&result, string(query), pq.Array(orderIDs))
	if err := validators.ValidateAssignedOrders(err, result); err != nil {
		return nil, err
	}
	query, _ = os.ReadFile(fmt.Sprintf("%s/%s.sql", folder, files[1]))
	err = tx.Get(&result1, string(query), pq.Array(courierIDs))
	if err := validators.ValidateExistingCouriers(err, result1); err != nil {
		return nil, err
	}
	chunkSize := 21845
	for i := 0; i < len(orders); i += chunkSize {
		end := i + chunkSize
		if end > len(orders) {
			end = len(orders)
		}
		chunk := orders[i:end]
		var placeholders []string
		var values []interface{}
		var orderIDs []int64
		for i, order := range chunk {
			placeholder := fmt.Sprintf("($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3)
			placeholders = append(placeholders, placeholder)
			values = append(values, order.OrderID, order.CourierId, order.CompleteTime)
		}
		firstQuery := "INSERT INTO order_completion (order_id, courier_id, complete_time) VALUES " + strings.Join(placeholders, ", ")
		if _, err := tx.Exec(firstQuery, values...); err != nil {
			return nil, err
		}
		for _, order := range chunk {
			orderIDs = append(orderIDs, order.OrderID)
		}
		query, _ = os.ReadFile(fmt.Sprintf("%s/%s.sql", folder, files[2]))
		rows, err := tx.Queryx(string(query), pq.Array(orderIDs))
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var completedOrder models.Order
			if err := rows.StructScan(&completedOrder); err != nil {
				return nil, err
			}
			completedOrders = append(completedOrders, completedOrder)
		}
		if err := rows.Err(); err != nil {
			return nil, err
		}
	}
	return completedOrders, nil
}
