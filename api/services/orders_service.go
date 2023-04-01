package services

import (
	"github.com/jmoiron/sqlx"
	"yandex-team.ru/bstask/api/models"
)

func GetOrders(db *sqlx.DB, limit int, offset int) ([]models.Order, error) {
	var orders []models.Order
	query := `SELECT id, cost, delivery_hours, delivery_district, weight FROM orders LIMIT $1 OFFSET $2`
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
