package validators

import (
	"fmt"
	"regexp"
	"yandex-team.ru/bstask/api/models"
)

type ValidationCourierError struct {
	Message string
	Data    models.CreateCourierDto
}

func (e *ValidationCourierError) Error() string {
	return e.Message
}

func ValidateTime(time string) bool {
	re := regexp.MustCompile(`^[0-9][0-9]:[0-9][0-9]-[0-9][0-9]:[0-9][0-9]$`)
	if re.MatchString(time) {
		var startHour, startMinute, endHour, endMinute int
		_, err := fmt.Sscanf(time, "%d:%d-%d:%d", &startHour, &startMinute, &endHour, &endMinute)
		if err != nil {
			return false
		}
		if startHour < 0 || startHour > 23 || endHour < 0 || endHour > 23 || startMinute < 0 || startMinute > 59 || endMinute < 0 || endMinute > 59 {
			return false
		} else {
			if startHour > endHour || (startHour == endHour && (startMinute > endMinute || startMinute == endMinute)) {
				return false
			} else {
				return true
			}
		}
	} else {
		return false
	}
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
		if courier.WorkingAreas[i] < 0 {
			return &ValidationCourierError{
				Message: "Courier working area is invalid",
				Data:    courier,
			}
		}
	}

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
	}
	return nil
}
