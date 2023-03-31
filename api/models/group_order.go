package models

type CouriersGroupOrders struct {
	CourierID      int64   `json:"courier_id"`
	AssignedOrders []int64 `json:"assigned_orders"`
	TotalCost      int     `json:"total_cost"`
}

type CreateGroupOrderDto struct {
	OrderIDs []int64 `json:"order_ids"`
}

type CreateGroupOrderRequest struct {
	Groups []CreateGroupOrderDto `json:"groups"`
}

type GroupOrderAssignResponse struct {
	Date   string                `json:"date"`
	Groups []CouriersGroupOrders `json:"groups"`
}
