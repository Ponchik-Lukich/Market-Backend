package services

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"os"
	"strings"
	"time"
	"yandex-team.ru/bstask/api/middleware"
	"yandex-team.ru/bstask/api/models"
	"yandex-team.ru/bstask/api/utils/validators"
)

func GetOrders(db *sqlx.DB, limit int, offset int) ([]models.Order, error) {
	var orders []models.Order
	query := `SELECT id, cost, delivery_hours, delivery_district, weight, complete_time FROM orders LIMIT $1 OFFSET $2`
	err := db.Select(&orders, query, limit, offset)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func GetOrderById(db *sqlx.DB, OrderID int64) (*models.Order, error) {
	var order models.Order
	query := `SELECT id, cost, delivery_hours, delivery_district, weight, complete_time FROM orders WHERE id = $1`
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
	var completeTime []*time.Time
	for _, order := range orders {
		if err := validators.ValidateCompleteOrder(order); err != nil {
			return nil, err
		}
		if err := validators.ValidateIds(order.OrderID, &setIDs); err != nil {
			return nil, err
		}
		orderIDs = append(orderIDs, order.OrderID)
		courierIDs = append(courierIDs, order.CourierId)
		completeTime = append(completeTime, order.CompleteTime)
	}
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	folder := fmt.Sprintf("%s/%s/%s/%s", dir, "api", "models", "queries")
	var ordersRes int
	var couriersRes []int64
	query, _ := os.ReadFile(fmt.Sprintf("%s/%s.sql", folder, "check_couriers"))
	err = tx.Select(&couriersRes, string(query), pq.Array(courierIDs))
	if err := validators.ValidateExistingCouriers(err, couriersRes); err != nil {
		return nil, err
	}
	query, _ = os.ReadFile(fmt.Sprintf("%s/%s.sql", folder, "check_orders"))
	err = tx.Get(&ordersRes, string(query), pq.Array(orderIDs), pq.Array(courierIDs))
	if err := validators.ValidateAssignedOrders(err, ordersRes); err != nil {
		return nil, err
	}
	query, _ = os.ReadFile(fmt.Sprintf("%s/%s.sql", folder, "update_orders"))
	rows, err := tx.Queryx(string(query), pq.Array(orderIDs), pq.Array(completeTime))
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
	return completedOrders, nil
}
