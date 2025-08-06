package handler

import (
	"net/http"

	"github.com/gabrielluizsf/rinha-backend-2005/db"
)

func PurgePayments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := db.Purge("payments"); err != nil {
		http.Error(w, "Failed to purge payments", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
