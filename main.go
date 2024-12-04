package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jontitorr/receipt-processor/handlers"
)

func main() {
	router := mux.NewRouter()
	store := handlers.NewReceiptStore()

	router.HandleFunc("/receipts/process", store.ProcessReceipt).Methods("POST")
	router.HandleFunc("/receipts/{id}/points", store.GetPoints).Methods("GET")

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
