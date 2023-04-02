package validators

import "yandex-team.ru/bstask/api/models"

type ValidationCourierError struct {
	Message string
	Data    models.CreateCourierDto
}

func (e *ValidationCourierError) Error() string {
	return e.Message
}

// TODO add validation

func ValidateCourier(courier models.CreateCourierDto) error {
	if courier.CourierType != models.FOOT && courier.CourierType != models.BIKE && courier.CourierType != models.AUTO {
		return &ValidationCourierError{
			Message: "Courier type is invalid",
			Data:    courier,
		}
	}
	//if len(courier.WorkingAreas) == 0 {
	//	return &ValidationCourierError{
	//		Message: "Courier working areas is empty",
	//		Data:    courier,
	//	}
	//}
	//if len(courier.WorkingHours) == 0 {
	//	return &ValidationCourierError{
	//		Message: "Courier working hours is empty",
	//		Data:    courier,
	//	}
	//}
	return nil
}
