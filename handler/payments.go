package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gabrielluizsf/rinha-backend-2005/adapter"
	"github.com/gabrielluizsf/rinha-backend-2005/db"
	"github.com/gabrielluizsf/rinha-backend-2005/types"
)

func Payments(c adapter.RequestContext) error {
	var p types.PaymentRequest
	if err := c.BodyParser(&p); err != nil {
		return c.Status(http.StatusBadRequest).SendString("Invalid JSON")
	}
	data, err := json.Marshal(p)
	if err != nil {
		statusCode := http.StatusInternalServerError
		statusText := http.StatusText(statusCode)
		return c.Status(statusCode).SendString(statusText)
	}
	if err := db.Save("payments_pending", data); err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Failed to queue payment")
	}

	return c.SendStatus(http.StatusAccepted)
}
