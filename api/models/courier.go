package models

type CourierType string

const (
	FOOT CourierType = "FOOT"
	BIKE CourierType = "BIKE"
	AUTO CourierType = "AUTO"
)

type Courier struct {
	CourierID    int64
	CourierType  CourierType
	Regions      []int
	WorkingHours []string
}

type CreateCourierDto struct {
	CourierType  CourierType `json:"courier_type"`
	Regions      []int       `json:"regions"`
	WorkingHours []string    `json:"working_hours"`
}

type CreateCourierRequest struct {
	Couriers []CreateCourierDto `json:"couriers"`
}

type CreateCouriersResponse struct {
	Couriers []Courier `json:"couriers"`
}

type GetCouriersResponse struct {
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
