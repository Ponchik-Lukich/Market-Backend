package services

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"os"
	"time"
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

func CreateCouriers(db *sqlx.DB, couriers []models.CreateCourierDto) ([]models.Courier, error) {
	var createdCouriers []models.Courier
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

func GetCourierMetaInfo(db *sqlx.DB, courierID int64, startDate string, endDate string) (*models.GetCourierMetaInfoResponse, error) {
	var courierMetaInfo models.GetCourierMetaInfoResponse
	if err := validators.ValidateCourierMeta(courierID, startDate, endDate); err != nil {
		return nil, err
	}
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	folder := fmt.Sprintf("%s/%s/%s/%s", dir, "api", "models", "queries")
	file := "return_rating"
	query, err := os.ReadFile(fmt.Sprintf("%s/%s.sql", folder, file))
	if err != nil {
		return nil, err
	}
	startDate = startDate + " 00:00:00"
	endDate = endDate + " 00:00:00"
	err = db.Get(&courierMetaInfo, string(query), courierID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	layout := "2006-01-02 15:04:05"
	start, _ := time.Parse(layout, startDate)
	end, _ := time.Parse(layout, endDate)
	duration := end.Sub(start)
	hours := int32(duration.Hours())
	switch courierMetaInfo.CourierType {
	case models.FOOT:
		courierMetaInfo.Earnings = courierMetaInfo.Earnings * 2
		courierMetaInfo.Rating = int32(float32(courierMetaInfo.Rating) / float32(hours) * 3)
	case models.BIKE:
		courierMetaInfo.Earnings = courierMetaInfo.Earnings * 3
		courierMetaInfo.Rating = int32(float32(courierMetaInfo.Rating) / float32(hours) * 2)
	case models.AUTO:
		courierMetaInfo.Earnings = courierMetaInfo.Earnings * 4
		courierMetaInfo.Rating = courierMetaInfo.Rating / hours
	}
	return &courierMetaInfo, nil
}
