package models

import (
	"github.com/lib/pq"
	"time"
)

type Order struct {
	OrderID       int64          `json:"order_id" db:"id"`
	Cost          int            `json:"cost" db:"cost"`
	DeliveryHours pq.StringArray `json:"delivery_hours" db:"delivery_hours"`
	Regions       int            `json:"regions" db:"delivery_district"`
	Weight        float64        `json:"weight" db:"weight"`
	CompletedTime *time.Time     `json:"completed_time,omitempty" db:"completion_time,omitempty"`
}

type CreateOrderDto struct {
	Cost          int32          `json:"cost" db:"cost"`
	Weight        float64        `json:"weight" db:"weight"`
	Regions       int32          `json:"regions" db:"delivery_district"`
	DeliveryHours pq.StringArray `json:"delivery_hours" db:"delivery_hours"`
}

type CreateOrderRequest struct {
	Orders []CreateOrderDto `json:"orders"`
}

type CreateOrderResponse struct {
	Orders []Order `json:"orders"`
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

type CompleteOrderRequestDto struct {
	Orders []CompleteOrder `json:"orders"`
}

type CompleteOrder struct {
	OrderID       int64      `json:"order_id" db:"id"`
	CourierId     int64      `json:"courier_id" db:"courier_id"`
	CompletedTime *time.Time `json:"completed_time" db:"completed_time"`
}
