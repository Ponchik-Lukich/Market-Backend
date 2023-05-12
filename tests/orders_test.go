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

	resp, err := http.Post("http://application:8080/orders", "application/json", bytes.NewBuffer(data))
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

func TestCreateOrderInvalidCost(t *testing.T) {
	order := models.CreateOrderDto{
		Cost:          0,
		DeliveryHours: []string{"10:00-12:00", "13:00-18:00"},
		Regions:       1,
		Weight:        1.0,
	}
	orderRequest := models.CreateOrderRequest{
		Orders: []models.CreateOrderDto{order},
	}

	data, err := json.Marshal(orderRequest)
	assert.NoError(t, err, "failed to marshal order")

	resp, err := http.Post("http://application:8080/orders", "application/json", bytes.NewBuffer(data))
	assert.NoError(t, err, "HTTP error")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "HTTP status code")
}

func TestCreateOrderInvalidWeight(t *testing.T) {
	order := models.CreateOrderDto{
		Cost:          1,
		DeliveryHours: []string{"10:00-12:00", "13:00-18:00"},
		Regions:       1,
		Weight:        0.0,
	}
	orderRequest := models.CreateOrderRequest{
		Orders: []models.CreateOrderDto{order},
	}

	data, err := json.Marshal(orderRequest)
	assert.NoError(t, err, "failed to marshal order")

	resp, err := http.Post("http://application:8080/orders", "application/json", bytes.NewBuffer(data))
	assert.NoError(t, err, "HTTP error")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "HTTP status code")
}

func TestCreateOrderInvalidDeliveryHours(t *testing.T) {
	nonexistentDeliveryHours := []string{
		"23:59-26:00",
	}

	invalidDeliveryHours1 := []string{
		"23:59-00:01",
	}

	invalidDeliveryHours2 := []string{
		"14:00-13:00",
	}

	overlappingDeliveryHours := []string{
		"11:00-13:00",
		"10:00-12:00",
	}

	testCases := []struct {
		name         string
		workingHours []string
	}{
		{"nonexistent_hours", nonexistentDeliveryHours},
		{"invalid_hours_1", invalidDeliveryHours1},
		{"invalid_hours_2", invalidDeliveryHours2},
		{"overlapping_hours", overlappingDeliveryHours},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			order := models.CreateOrderDto{
				Cost:          1,
				DeliveryHours: tc.workingHours,
				Regions:       1,
				Weight:        1.0,
			}
			orderRequest := models.CreateOrderRequest{
				Orders: []models.CreateOrderDto{order},
			}

			data, err := json.Marshal(orderRequest)
			assert.NoError(t, err, "failed to marshal order")

			resp, err := http.Post("http://application:8080/orders", "application/json", bytes.NewBuffer(data))
			assert.NoError(t, err, "HTTP error")
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "HTTP status code")
		})
	}
}

func TestCreateOrderInvalidRegions(t *testing.T) {
	order := models.CreateOrderDto{
		Cost:          1,
		DeliveryHours: []string{"10:00-12:00", "13:00-18:00"},
		Regions:       0,
		Weight:        1.0,
	}
	orderRequest := models.CreateOrderRequest{
		Orders: []models.CreateOrderDto{order},
	}

	data, err := json.Marshal(orderRequest)
	assert.NoError(t, err, "failed to marshal order")

	resp, err := http.Post("http://application:8080/orders", "application/json", bytes.NewBuffer(data))
	assert.NoError(t, err, "HTTP error")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "HTTP status code")
}

func TestGetOrderById(t *testing.T) {
	var orderId int64 = 1

	resp, err := http.Get(fmt.Sprintf("http://application:8080/orders/%d", orderId))
	assert.NoError(t, err, "HTTP error")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "HTTP status code")

	var retrievedOrder models.Order
	err = json.NewDecoder(resp.Body).Decode(&retrievedOrder)
	assert.NoError(t, err, "failed to decode HTTP body")

	assert.Equal(t, orderId, retrievedOrder.OrderID, "Mismatch in order ID")
}

func TestGetOrderByIdInvalid(t *testing.T) {
	var orderId int64 = -1

	resp, err := http.Get(fmt.Sprintf("http://application:8080/orders/%d", orderId))
	assert.NoError(t, err, "HTTP error")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "HTTP status code")
}
