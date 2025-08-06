package routes

import (
	"github.com/gabrielluizsf/rinha-backend-2005/adapter"
	"github.com/gabrielluizsf/rinha-backend-2005/handler"
)

func InitRoutes(app adapter.ServerManager) {
	app.Post("/payments", handler.Payments)
	app.Get("/payments-summary", handler.PaymentsSummary)
	app.Post("/purge-payments", handler.PurgePayments)
}
