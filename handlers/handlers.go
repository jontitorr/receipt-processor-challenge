package handlers

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jontitorr/receipt-processor/models"
	"github.com/jontitorr/receipt-processor/service"
)

type ReceiptStore struct {
	receipts map[string]*models.Receipt
	mu       sync.RWMutex
}

func NewReceiptStore() *ReceiptStore {
	return &ReceiptStore{
		receipts: make(map[string]*models.Receipt),
	}
}

func (rs *ReceiptStore) ProcessReceipt(w http.ResponseWriter, r *http.Request) {
	var receipt models.Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, "Invalid receipt format", http.StatusBadRequest)
		return
	}

	// Basic validation
	if receipt.Retailer == "" || receipt.PurchaseDate == "" || receipt.PurchaseTime == "" || 
		receipt.Total == "" || len(receipt.Items) == 0 {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	id := uuid.New().String()
	rs.mu.Lock()
	rs.receipts[id] = &receipt
	rs.mu.Unlock()

	response := models.ProcessResponse{ID: id}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (rs *ReceiptStore) GetPoints(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	rs.mu.RLock()
	receipt, exists := rs.receipts[id]
	rs.mu.RUnlock()

	if !exists {
		http.Error(w, "No receipt found for that id", http.StatusNotFound)
		return
	}

	points := service.CalculatePoints(receipt)
	response := models.PointsResponse{Points: points}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
