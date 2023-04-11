package validators

import (
	"yandex-team.ru/bstask/api/models"
)

type ValidationCourierError struct {
	Message string
	Data    models.CreateCourierDto
}

func (e *ValidationCourierError) Error() string {
	return e.Message
}

func ValidateCourier(courier models.CreateCourierDto) error {
	if courier.CourierType != models.FOOT && courier.CourierType != models.BIKE && courier.CourierType != models.AUTO {
		return &ValidationCourierError{
			Message: "Courier type is invalid",
			Data:    courier,
		}
	}
	setAreas := map[int64]struct{}{}
	for i := 0; i < len(courier.WorkingAreas); i++ {
		if _, ok := setAreas[courier.WorkingAreas[i]]; ok {
			return &ValidationCourierError{
				Message: "Courier working areas has duplicates",
				Data:    courier,
			}
		} else {
			setAreas[courier.WorkingAreas[i]] = struct{}{}
		}
		if courier.WorkingAreas[i] <= 0 {
			return &ValidationCourierError{
				Message: "Courier working area is invalid",
				Data:    courier,
			}
		}
	}
	var intervals []string
	setHours := map[string]struct{}{}
	for i := 0; i < len(courier.WorkingHours); i++ {
		if _, ok := setHours[courier.WorkingHours[i]]; ok {
			return &ValidationCourierError{
				Message: "Courier working hours has duplicates",
				Data:    courier,
			}
		} else {
			setHours[courier.WorkingHours[i]] = struct{}{}
		}
		if !ValidateTime(courier.WorkingHours[i]) {
			return &ValidationCourierError{
				Message: "Courier working hour is invalid",
				Data:    courier,
			}
		}
		intervals = append(intervals, courier.WorkingHours[i])
	}
	if !ValidateTimeIntervals(intervals) {
		return &ValidationCourierError{
			Message: "Courier working hours are overlapping",
			Data:    courier,
		}
	}
	return nil
}
