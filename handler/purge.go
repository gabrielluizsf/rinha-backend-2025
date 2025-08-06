package handler

import (
	"net/http"

	"github.com/gabrielluizsf/rinha-backend-2005/adapter"
	"github.com/gabrielluizsf/rinha-backend-2005/db"
)

func PurgePayments(c adapter.RequestContext) error {
	if err := db.Purge("payments"); err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Failed to purge payments")
	}
	return c.SendStatus(http.StatusOK)
}
