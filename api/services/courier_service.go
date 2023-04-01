package services

import (
	"github.com/jmoiron/sqlx"
	"yandex-team.ru/bstask/api/models"
)

//func GetCourier(db *sqlx.DB, courierID int64) (*models.Courier, error) {
// Implement logic for getting a courier from the database
//query := `SELECT * FROM couriers WHERE id = $1`
//return nil, nil
//}

func GetCouriers(db *sqlx.DB, limit int, offset int) ([]models.Courier, error) {
	//Implement logic for getting couriers from the database
	println("In GetCouriers")
	var couriers []models.Courier
	query := `SELECT * FROM couriers LIMIT $1 OFFSET $2`
	rows, err := db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var courier models.Courier
		err := rows.Scan(&courier.CourierID, &courier.CourierType, &courier.MaxWeight)
		if err != nil {
			return nil, err
		}
		couriers = append(couriers, courier)
	}
	return couriers, nil
}

//func CreateCourier(db *sqlx.DB, courier *models.Courier) error {
//Implement logic for creating a courier in the database
//query := `INSERT INTO couriers (id, name, phone, working_hours) VALUES ($1, $2, $3, $4)`
//return nil
//}

//func UpdateCourier(db *sqlx.DB, courier *models.Courier) error {
// Implement logic for updating a courier in the database
//query := `UPDATE couriers SET name = $1, phone = $2, working_hours = $3 WHERE id = $4`
//return nil
//}
