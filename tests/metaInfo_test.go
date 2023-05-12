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

func TestGetCourierMetaIfo(t *testing.T) {
	var orders []models.CreateOrderDto
	costs := []int32{100, 200, 300, 100, 200, 300, 100, 200, 300}
	deliveryHours := "00:00-23:59"
	for i := 0; i < 9; i++ {
		order := models.CreateOrderDto{
			Cost:          costs[i],
			DeliveryHours: []string{deliveryHours},
			Regions:       1,
			Weight:        0.01,
		}
		orders = append(orders, order)
	}
	orderRequest := models.CreateOrderRequest{
		Orders: orders,
	}

	data, err := json.Marshal(orderRequest)
	assert.NoError(t, err, "failed to marshal order")

	resp, err := http.Post("http://app:8080/orders", "application/json", bytes.NewBuffer(data))
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
	if err = json.Unmarshal(body, &response); err != nil {
		t.Fatal(err)
	}

	if len(response.Orders) == 0 {
		t.Fatal("response did not contain any orders")
	}
	var createdOrderIds []int64
	for i := 0; i < len(response.Orders); i++ {
		createdOrderIds = append(createdOrderIds, response.Orders[i].OrderID)
	}

	var couriers []models.CreateCourierDto
	couriers = append(couriers, models.CreateCourierDto{
		WorkingHours: []string{"00:00-23:59"},
		WorkingAreas: []int64{1},
		CourierType:  "FOOT",
	})
	courierRequest := models.CreateCourierRequest{
		Couriers: couriers,
	}

	data, err = json.Marshal(courierRequest)
	assert.NoError(t, err, "failed to marshal couriers")

	resp, err = http.Post("http://app:8080/couriers", "application/json", bytes.NewBuffer(data))
	assert.NoError(t, err, "HTTP error")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "HTTP status code")

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var courierResponse struct {
		Couriers []models.Courier `json:"couriers"`
	}
	if err = json.Unmarshal(body, &courierResponse); err != nil {
		t.Fatal(err)
	}

	if len(courierResponse.Couriers) == 0 {
		t.Fatal("response did not contain any couriers")
	}
	courierID := courierResponse.Couriers[0].CourierID
	courierWorkingHours := courierResponse.Couriers[0].WorkingHours
	courierWorkingAreas := courierResponse.Couriers[0].WorkingAreas
	courierType := courierResponse.Couriers[0].CourierType
	resp, err = http.Get(fmt.Sprintf("http://app:8080/test/%d", courierID))
	assert.NoError(t, err, "HTTP error")
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode, "HTTP status code")
	var completedOrders []models.CompleteOrderDto
	var completeTime *time.Time
	ti := time.Date(2023, 1, 1, 6, 0, 0, 0, time.UTC)
	completeTime = &ti

	for i := 0; i < len(createdOrderIds); i++ {
		completedOrder := models.CompleteOrderDto{
			OrderID:      createdOrderIds[i],
			CourierId:    courierID,
			CompleteTime: completeTime,
		}
		completedOrders = append(completedOrders, completedOrder)
	}

	completeOrderRequest := models.CompleteOrderRequest{
		Orders: completedOrders,
	}

	data, err = json.Marshal(completeOrderRequest)
	assert.NoError(t, err, "failed to marshal complete order")

	resp, err = http.Post("http://app:8080/orders/complete", "application/json", bytes.NewBuffer(data))
	assert.NoError(t, err, "HTTP error")
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode, "HTTP status code")

	startDate := "2023-01-01"
	endDate := "2023-01-02"
	resp, err = http.Get(fmt.Sprintf("http://app:8080/couriers/meta-info/%d?startDate=%s&endDate=%s", courierID, startDate, endDate))
	assert.NoError(t, err, "HTTP error")
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode, "HTTP status code")

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	result := models.GetCourierMetaInfoResponse{}
	if err = json.Unmarshal(body, &result); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, courierID, result.CourierID, "CourierID")
	assert.Equal(t, courierType, result.CourierType, "Type")
	assert.Equal(t, courierWorkingHours, result.WorkingHours, "WorkingHours")
	assert.Equal(t, courierWorkingAreas, result.WorkingAreas, "WorkingAreas")
	assert.Equal(t, int32(1), result.Rating, "Rating")
	assert.Equal(t, int32(3600), result.Earnings, "Earnings")
}
