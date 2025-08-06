package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gabrielluizsf/rinha-backend-2005/db"
	"github.com/gabrielluizsf/rinha-backend-2005/types"
)

func Payments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var p types.PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(p)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := db.Save("payments_pending", data); err != nil {
		http.Error(w, "Failed to queue payment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
