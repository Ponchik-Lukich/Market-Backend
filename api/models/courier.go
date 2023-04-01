package models

import "github.com/lib/pq"

type CourierType string

const (
	FOOT CourierType = "FOOT"
	BIKE CourierType = "BIKE"
	AUTO CourierType = "AUTO"
)

type Courier struct {
	CourierID    int64          `json:"courier_id" db:"id"`
	CourierType  CourierType    `json:"courier_type" db:"type" `
	WorkingAreas pq.Int64Array  `json:"regions" db:"working_areas"`
	WorkingHours pq.StringArray `json:"working_hours" db:"working_hours"`
}

type CreateCourierDto struct {
	CourierType  CourierType    `json:"courier_type" db:"type" `
	WorkingAreas pq.Int64Array  `json:"regions" db:"working_areas"`
	WorkingHours pq.StringArray `json:"working_hours" db:"working_hours"`
}

type CreateCourierRequest struct {
	Couriers []CreateCourierDto `json:"couriers"`
}

type CreateCourierResponse struct {
	Couriers []Courier `json:"couriers"`
}

type GetCourierResponse struct {
	Couriers []Courier `json:"couriers"`
	Limit    int       `json:"limit"`
	Offset   int       `json:"offset"`
}

type GetCourierMetaInfoResponse struct {
	CourierID    int64       `json:"courier_id"`
	CourierType  CourierType `json:"courier_type"`
	Regions      []int       `json:"regions"`
	WorkingHours []string    `json:"working_hours"`
	Rating       int         `json:"rating"`
	Earnings     int         `json:"earnings"`
}
