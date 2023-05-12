package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
	"time"
	"yandex-team.ru/bstask/api/models"
)

func waitForAppReady() {
	maxRetries := 10
	retries := 0
	for retries < maxRetries {
		resp, err := http.Get("http://application:8080/ping")
		if err == nil && resp.StatusCode == http.StatusOK {
			break
		}
		retries++
		time.Sleep(1 * time.Second)
	}
}

func TestCreateCourier(t *testing.T) {
	waitForAppReady()
	courier := models.CreateCourierDto{
		WorkingHours: []string{"10:00-12:00", "13:00-18:00"},
		WorkingAreas: []int64{1, 2, 3},
		CourierType:  "FOOT",
	}
	courierRequest := models.CreateCourierRequest{
		Couriers: []models.CreateCourierDto{courier},
	}

	data, err := json.Marshal(courierRequest)
	assert.NoError(t, err, "failed to marshal courier")

	resp, err := http.Post("http://application:8080/couriers", "application/json", bytes.NewBuffer(data))
	assert.NoError(t, err, "HTTP error")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "HTTP status code")

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var response struct {
		Couriers []models.Courier `json:"couriers"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatal(err)
	}

	if len(response.Couriers) == 0 {
		t.Fatal("response did not contain any couriers")
	}

	createdCourier := response.Couriers[0]
	assert.Equal(t, courier.WorkingHours, createdCourier.WorkingHours, "Mismatch in courier working hours")
	assert.Equal(t, courier.WorkingAreas, createdCourier.WorkingAreas, "Mismatch in courier working areas")
	assert.Equal(t, courier.CourierType, createdCourier.CourierType, "Mismatch in courier type")
}

func TestCreateCourierInvalidWorkingHours(t *testing.T) {
	nonexistentWorkingHours := []string{
		"25:00-26:00",
	}

	invalidWorkingHours1 := []string{
		"23:59-00:01",
	}

	invalidWorkingHours2 := []string{
		"14:00-13:00",
	}

	overlappingWorkingHours := []string{
		"11:00-13:00",
		"10:00-12:00",
	}

	testCases := []struct {
		name         string
		workingHours []string
	}{
		{"nonexistent_hours", nonexistentWorkingHours},
		{"invalid_hours_1", invalidWorkingHours1},
		{"invalid_hours_2", invalidWorkingHours2},
		{"overlapping_hours", overlappingWorkingHours},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			courier := models.CreateCourierDto{
				WorkingHours: tc.workingHours,
				WorkingAreas: []int64{1, 2, 3},
				CourierType:  "FOOT",
			}
			courierRequest := models.CreateCourierRequest{
				Couriers: []models.CreateCourierDto{courier},
			}

			data, err := json.Marshal(courierRequest)
			assert.NoError(t, err, "failed to marshal courier")

			resp, err := http.Post("http://application:8080/couriers", "application/json", bytes.NewBuffer(data))
			assert.NoError(t, err, "HTTP error")
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "HTTP status code")
		})
	}
}

func TestCreateCourierInvalidWorkingAreas(t *testing.T) {

	zeroAreas := []int64{0}

	negativeAreas := []int64{1, -1}

	overlappingAreas := []int64{1, 2, 1}

	testCases := []struct {
		name         string
		workingAreas []int64
	}{
		{"zero_areas", zeroAreas},
		{"negative_areas", negativeAreas},
		{"overlapping_areas", overlappingAreas},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			courier := models.CreateCourierDto{
				WorkingHours: []string{"10:00-12:00"},
				WorkingAreas: tc.workingAreas,
				CourierType:  "FOOT",
			}
			courierRequest := models.CreateCourierRequest{
				Couriers: []models.CreateCourierDto{courier},
			}

			data, err := json.Marshal(courierRequest)
			assert.NoError(t, err, "failed to marshal courier")

			resp, err := http.Post("http://application:8080/couriers", "application/json", bytes.NewBuffer(data))
			assert.NoError(t, err, "HTTP error")
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "HTTP status code")
		})
	}
}

func TestCreateCourierInvalidType(t *testing.T) {
	courier := models.CreateCourierDto{
		WorkingHours: []string{"10:00-12:00", "13:00-18:00"},
		WorkingAreas: []int64{1, 2, 3},
		CourierType:  "FOoT",
	}
	courierRequest := models.CreateCourierRequest{
		Couriers: []models.CreateCourierDto{courier},
	}

	data, err := json.Marshal(courierRequest)
	assert.NoError(t, err, "failed to marshal courier")

	resp, err := http.Post("http://application:8080/couriers", "application/json", bytes.NewBuffer(data))
	assert.NoError(t, err, "HTTP error")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "HTTP status code")
}

func TestGetCourierById(t *testing.T) {
	var courierID int64 = 1
	resp, err := http.Get(fmt.Sprintf("http://application:8080/couriers/%d", courierID))
	assert.NoError(t, err, "HTTP error")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "HTTP status code")

	var retrievedCourier models.Courier
	if err := json.NewDecoder(resp.Body).Decode(&retrievedCourier); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, courierID, retrievedCourier.CourierID, "Mismatch in courier ID")
}

func TestGetCourierByIdInvalid(t *testing.T) {
	var courierID int64 = -1
	resp, err := http.Get(fmt.Sprintf("http://application:8080/couriers/%d", courierID))
	assert.NoError(t, err, "HTTP error")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "HTTP status code")
}
