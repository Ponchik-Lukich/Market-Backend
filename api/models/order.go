package models

import (
	"github.com/lib/pq"
)

type Order struct {
	OrderID       int64          `json:"order_id" db:"id"`
	Cost          int32          `json:"cost" db:"cost"`
	DeliveryHours pq.StringArray `json:"delivery_hours" db:"delivery_hours"`
	Regions       int32          `json:"regions" db:"delivery_district"`
	Weight        float64        `json:"weight" db:"weight"`
	//CompletedTime *time.Time     `json:"completed_time" db:"completed_time"`
}

type CreateOrderDto struct {
	Cost          int32          `json:"cost" db:"cost"`
	DeliveryHours pq.StringArray `json:"delivery_hours" db:"delivery_hours"`
	Regions       int32          `json:"regions" db:"delivery_district"`
	Weight        float64        `json:"weight" db:"weight"`
}

type CreateOrderRequest struct {
	Orders []CreateOrderDto `json:"orders"`
}

type GetOrderResponse struct {
	Orders []Order `json:"orders"`
	Limit  int     `json:"limit"`
	Offset int     `json:"offset"`
}

type OrderAssignResponse struct {
	Date     string                `json:"date"`
	Couriers []CouriersGroupOrders `json:"couriers"`
}
