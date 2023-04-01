package validators

import "yandex-team.ru/bstask/api/models"

type ValidationError struct {
	Message string
	Data    models.CreateCourierDto
}

func (e *ValidationError) Error() string {
	return e.Message
}

func ValidateCourier(courier models.CreateCourierDto) error {
	if courier.CourierType != models.FOOT && courier.CourierType != models.BIKE && courier.CourierType != models.AUTO {
		return &ValidationError{
			Message: "Courier type is invalid",
			Data:    courier,
		}
	}
	if len(courier.WorkingAreas) == 0 {
		return &ValidationError{
			Message: "Courier working areas is empty",
			Data:    courier,
		}
	}
	if len(courier.WorkingHours) == 0 {
		return &ValidationError{
			Message: "Courier working hours is empty",
			Data:    courier,
		}
	}
	return nil
}
