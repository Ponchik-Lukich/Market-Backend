package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
	"yandex-team.ru/bstask/api/models"
)

func TestCreateCourier(t *testing.T) {
	// Create a new courier
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

	resp, err := http.Post("http://localhost:8080/couriers", "application/json", bytes.NewBuffer(data))
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

			resp, err := http.Post("http://localhost:8080/couriers", "application/json", bytes.NewBuffer(data))
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

			resp, err := http.Post("http://localhost:8080/couriers", "application/json", bytes.NewBuffer(data))
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

	resp, err := http.Post("http://localhost:8080/couriers", "application/json", bytes.NewBuffer(data))
	assert.NoError(t, err, "HTTP error")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "HTTP status code")
}

func TestGetCourierById(t *testing.T) {
	var courierID int64 = 1
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/couriers/%d", courierID))
	assert.NoError(t, err, "HTTP error")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "HTTP status code")

	var retrievedCourier models.Courier
	if err := json.NewDecoder(resp.Body).Decode(&retrievedCourier); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, courierID, retrievedCourier.CourierID, "Mismatch in courier ID")
}

func TestCreateOrder(t *testing.T) {
	order := models.CreateOrderDto{
		Cost:          1,
		DeliveryHours: []string{"10:00-12:00", "13:00-18:00"},
		Regions:       1,
		Weight:        1.0,
	}
	orderRequest := models.CreateOrderRequest{
		Orders: []models.CreateOrderDto{order},
	}

	data, err := json.Marshal(orderRequest)
	assert.NoError(t, err, "failed to marshal order")

	resp, err := http.Post("http://localhost:8080/orders", "application/json", bytes.NewBuffer(data))
	assert.NoError(t, err, "HTTP error")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "HTTP status code")

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var response struct {
		Orders []models.Order `json:"orders"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatal(err)
	}

	if len(response.Orders) == 0 {
		t.Fatal("response did not contain any orders")
	}

	createdOrder := response.Orders[0]
	assert.Equal(t, order.Cost, createdOrder.Cost, "Mismatch in order cost")
	assert.Equal(t, order.DeliveryHours, createdOrder.DeliveryHours, "Mismatch in order delivery hours")
	assert.Equal(t, order.Regions, createdOrder.Regions, "Mismatch in order regions")
	assert.Equal(t, order.Weight, createdOrder.Weight, "Mismatch in order weight")

}

func TestGetOrder(t *testing.T) {
	var orderId int64 = 1

	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/orders/%d", orderId))
	assert.NoError(t, err, "HTTP error")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "HTTP status code")

	var retrievedOrder models.Order
	err = json.NewDecoder(resp.Body).Decode(&retrievedOrder)
	assert.NoError(t, err, "failed to decode HTTP body")

	assert.Equal(t, orderId, retrievedOrder.OrderID, "Mismatch in order ID")
}
