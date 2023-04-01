package services

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"net/http"
	"yandex-team.ru/bstask/api/models"
	"yandex-team.ru/bstask/api/utils/validators"
)

func GetCourierById(db *sqlx.DB, courierID int64) (*models.Courier, error) {
	var courier models.Courier
	query := `SELECT * FROM couriers WHERE id = $1`
	err := db.Get(&courier, query, courierID)
	if err != nil {
		return nil, err
	}
	return &courier, nil
}

func GetCouriers(db *sqlx.DB, limit int, offset int) ([]models.Courier, error) {
	var couriers []models.Courier
	query := `SELECT * FROM couriers LIMIT $1 OFFSET $2`
	err := db.Select(&couriers, query, limit, offset)

	if err != nil {
		panic(err)
		return nil, err
	}
	return couriers, nil
}

//func CreateCouriers(db *sqlx.DB, couriers []models.CreateCourierDto) error {
//
//	query := `INSERT INTO couriers (type, working_areas, working_hours) VALUES ($1, $2, $3)`
//	stmt, err := db.Prepare(query)
//	if err != nil {
//		return echo.NewHTTPError(http.StatusInternalServerError, "Error preparing query")
//	}
//	defer stmt.Close()
//
//	for _, courier := range couriers {
//		err = validators.ValidateCourier(courier)
//		if err != nil {
//			return &validators.ValidationError{
//				Message: "Validation failed for courier",
//				Data:    courier,
//			}
//		}
//		_, err := stmt.Exec(courier.CourierType, courier.WorkingAreas, courier.WorkingHours)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}

func CreateCouriers(db *sqlx.DB, couriers []models.CreateCourierDto) ([]models.Courier, error) {
	var createdCouriers []models.Courier
	query := `INSERT INTO couriers (type, working_areas, working_hours) VALUES ($1, $2, $3) RETURNING id, type, working_areas, working_hours`
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Error preparing query")
	}
	defer stmt.Close()

	for _, courier := range couriers {
		err = validators.ValidateCourier(courier)
		if err != nil {
			return nil, &validators.ValidationError{
				Message: "Validation failed for courier",
				Data:    courier,
			}
		}
		var createdCourier models.Courier
		err := stmt.QueryRow(courier.CourierType, courier.WorkingAreas, courier.WorkingHours).Scan(&createdCourier.CourierID, &createdCourier.CourierType, &createdCourier.WorkingAreas, &createdCourier.WorkingHours)
		if err != nil {
			return nil, err
		}
		createdCouriers = append(createdCouriers, createdCourier)
	}
	return createdCouriers, nil
}
