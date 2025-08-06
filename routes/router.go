package routes

import (
	"net/http"

	"github.com/gabrielluizsf/rinha-backend-2005/handler"
)

type RouterRegister interface {
	HandleFunc(endpoint string, handler func(http.ResponseWriter, *http.Request))
}

func InitRoutes(r RouterRegister) {
	r.HandleFunc("/payments", handler.Payments)
	r.HandleFunc("/payments-summary", handler.PaymentsSummary)
	r.HandleFunc("/purge-payments", handler.PurgePayments)
}
