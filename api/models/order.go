package models

import "time"

type Order struct {
	OrderID       int64
	Weight        float64
	Regions       int
	DeliveryHours []string
	Cost          int
	CompletedTime *time.Time
}

type CreateOrderDto struct {
	Weight        float64  `json:"weight"`
	Regions       int      `json:"regions"`
	DeliveryHours []string `json:"delivery_hours"`
	Cost          int      `json:"cost"`
}

type CreateOrderRequest struct {
	Orders []CreateOrderDto `json:"orders"`
}

type OrderAssignResponse struct {
	Date     string                `json:"date"`
	Couriers []CouriersGroupOrders `json:"couriers"`
}
