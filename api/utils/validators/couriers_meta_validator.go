package validators

import (
	"time"
)

type ValidationCourierMetaError struct {
	Message string
}

func (e *ValidationCourierMetaError) Error() string {
	return e.Message
}

func ValidateCourierMeta(courierID int64, startDate string, endDate string) error {
	if courierID <= 0 {
		return &ValidationCourierMetaError{
			Message: "Courier id is invalid",
		}
	}
	if startDate == "" || endDate == "" {
		return &ValidationCourierMetaError{
			Message: "startDate and endDate must be specified",
		}
	}
	layout := "2006-01-02"
	start, err := time.Parse(layout, startDate)
	if err != nil {
		return &ValidationCourierMetaError{
			Message: "startDate is invalid",
		}
	}
	end, err := time.Parse(layout, endDate)
	if err != nil {
		return &ValidationCourierMetaError{
			Message: "endDate is invalid",
		}
	}
	if start.After(end) || start.Equal(end) {
		return &ValidationCourierMetaError{
			Message: "startDate must be before endDate",
		}
	}
	return nil
}
