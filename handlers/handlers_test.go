package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jontitorr/receipt-processor/models"
)

func TestProcessReceipt(t *testing.T) {
	store := NewReceiptStore()

	tests := []struct {
		name           string
		receipt       models.Receipt
		expectedStatus int
	}{
		{
			name: "Valid Receipt",
			receipt: models.Receipt{
				Retailer:     "Target",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Items: []models.Item{
					{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
				},
				Total: "6.49",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid Receipt - Missing Required Fields",
			receipt: models.Receipt{
				Retailer: "Target",
				// Missing other required fields
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.receipt)
			req := httptest.NewRequest("POST", "/receipts/process", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			store.ProcessReceipt(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("ProcessReceipt() status = %v, want %v", w.Code, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusOK {
				var response models.ProcessResponse
				json.NewDecoder(w.Body).Decode(&response)
				if response.ID == "" {
					t.Error("ProcessReceipt() did not return an ID")
				}
			}
		})
	}
}

func TestGetPoints(t *testing.T) {
	store := NewReceiptStore()

	// Add a test receipt
	receipt := &models.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []models.Item{
			{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
		},
		Total: "6.49",
	}

	store.mu.Lock()
	store.receipts["test-id"] = receipt
	store.mu.Unlock()

	tests := []struct {
		name           string
		receiptID     string
		expectedStatus int
		expectPoints   bool
	}{
		{
			name:           "Valid Receipt ID",
			receiptID:     "test-id",
			expectedStatus: http.StatusOK,
			expectPoints:   true,
		},
		{
			name:           "Invalid Receipt ID",
			receiptID:     "non-existent-id",
			expectedStatus: http.StatusNotFound,
			expectPoints:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/receipts/"+tt.receiptID+"/points", nil)
			w := httptest.NewRecorder()

			// Setup router with path parameters
			router := mux.NewRouter()
			router.HandleFunc("/receipts/{id}/points", store.GetPoints)
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("GetPoints() status = %v, want %v", w.Code, tt.expectedStatus)
			}

			if tt.expectPoints {
				var response models.PointsResponse
				json.NewDecoder(w.Body).Decode(&response)
				if response.Points == 0 {
					t.Error("GetPoints() returned 0 points for valid receipt")
				}
			}
		})
	}
}
