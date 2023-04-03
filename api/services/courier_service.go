package services

import (
	"github.com/jmoiron/sqlx"
	"yandex-team.ru/bstask/api/models"
	"yandex-team.ru/bstask/api/utils/validators"
)

func GetCourierById(db *sqlx.DB, courierID int64) (*models.Courier, error) {
	var courier models.Courier
	query := `SELECT * FROM couriers WHERE id = $1`
	err := db.Get(&courier, query, courierID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, err
	}
	return &courier, nil
}

func GetCouriers(db *sqlx.DB, limit int, offset int) ([]models.Courier, error) {
	var couriers []models.Courier
	query := `SELECT * FROM couriers LIMIT $1 OFFSET $2`
	err := db.Select(&couriers, query, limit, offset)

	if err != nil {
		return nil, err
	}
	return couriers, nil
}

//
//func CreateCouriers(db *sqlx.DB, couriers []models.CreateCourierDto) ([]models.Courier, error) {
//	var createdCouriers []models.Courier
//
//	// Validate couriers
//	for _, courier := range couriers {
//		err := validators.ValidateCourier(courier)
//		if err != nil {
//			return nil, &validators.ValidationCourierError{
//				Message: "Validation failed for courier",
//				Data:    courier,
//			}
//		}
//	}
//
//	// Split the couriers into chunks of size 30000
//	chunkSize := 21500
//	chunks := make([][]models.CreateCourierDto, 0)
//	for i := 0; i < len(couriers); i += chunkSize {
//		end := i + chunkSize
//		if end > len(couriers) {
//			end = len(couriers)
//		}
//		chunks = append(chunks, couriers[i:end])
//	}
//
//	// Execute a separate insert statement for each chunk of couriers
//	for _, chunk := range chunks {
//		var query strings.Builder
//		query.WriteString("INSERT INTO couriers (type, working_areas, working_hours) VALUES ")
//
//		var placeholders []string
//		var values []interface{}
//		for i, courier := range chunk {
//			placeholder := fmt.Sprintf("($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3)
//			placeholders = append(placeholders, placeholder)
//			values = append(values, courier.CourierType, courier.WorkingAreas, courier.WorkingHours)
//		}
//
//		query.WriteString(strings.Join(placeholders, ", "))
//		query.WriteString(" RETURNING id, type, working_areas, working_hours")
//
//		rows, err := db.Query(query.String(), values...)
//		if err != nil {
//			return nil, err
//		}
//		defer rows.Close()
//
//		for rows.Next() {
//			var createdCourier models.Courier
//			err := rows.Scan(&createdCourier.CourierID, &createdCourier.CourierType, &createdCourier.WorkingAreas, &createdCourier.WorkingHours)
//			if err != nil {
//				return nil, err
//			}
//			createdCouriers = append(createdCouriers, createdCourier)
//		}
//
//		if err := rows.Err(); err != nil {
//			return nil, err
//		}
//	}
//
//	return createdCouriers, nil
//}

func CreateCouriers(db *sqlx.DB, couriers []models.CreateCourierDto) ([]models.Courier, error) {
	var createdCouriers []models.Courier
	// Validate couriers
	for _, courier := range couriers {
		if err := validators.ValidateCourier(courier); err != nil {
			return nil, err
		}
	}

	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	chunkSize := 21845
	for i := 0; i < len(couriers); i += chunkSize {
		end := i + chunkSize
		if end > len(couriers) {
			end = len(couriers)
		}
		chunk := couriers[i:end]

		query := `INSERT INTO couriers (type, working_areas, working_hours) VALUES (:type, :working_areas, :working_hours) RETURNING id, type, working_areas, working_hours`
		rows, err := tx.NamedQuery(query, chunk)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var createdCourier models.Courier
			if err := rows.StructScan(&createdCourier); err != nil {
				return nil, err
			}
			createdCouriers = append(createdCouriers, createdCourier)
		}

		if err := rows.Err(); err != nil {
			return nil, err
		}
	}

	return createdCouriers, nil
}
