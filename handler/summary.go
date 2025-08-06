package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gabrielluizsf/rinha-backend-2005/adapter"
	"github.com/gabrielluizsf/rinha-backend-2005/date"
	"github.com/gabrielluizsf/rinha-backend-2005/db"
	"github.com/gabrielluizsf/rinha-backend-2005/types"
	"github.com/i9si-sistemas/stringx"
)

func PaymentsSummary(c adapter.RequestContext) error {
	from, to, useFilter, err := parseTimeRange(c)
	if err != nil {
		if err == http.ErrMissingFile {
			return c.Status(http.StatusBadRequest).SendString("Both 'from' and 'to' query parameters are required, or omit both for all data")
		} else {
			return c.Status(http.StatusBadRequest).SendString("Invalid datetime format")
		}
	}
	paymentsData, err := db.GetAll("payments")
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Failed to retrieve payments from Redis")
	}
	resp := summarizePayments(paymentsData, from, to, useFilter)

	return c.JSON(resp)
}

func summarizePayments(
	paymentsData db.GetAllResult,
	from, to time.Time,
	useFilter bool,
) types.PaymentsSummaryResponse {
	var (
		defaultCount,
		fallbackCount int
	)
	var (
		defaultSum, fallbackSum float64
	)
	for _, paymentDataStr := range paymentsData {
		var payment types.ProcessedPayment
		if err := json.Unmarshal([]byte(paymentDataStr), &payment); err != nil {
			continue
		}
		if useFilter {
			createdAt, err := date.Parse(payment.CreatedAt)
			if err != nil || createdAt.Before(from) || createdAt.After(to) {
				continue
			}
		}
		if payment.Processor == "DEFAULT" && payment.Status == "PROCESSED_DEFAULT" {
			defaultCount++
			defaultSum += payment.Amount
		} else if payment.Processor == "FALLBACK" && payment.Status == "PROCESSED_FALLBACK" {
			fallbackCount++
			fallbackSum += payment.Amount
		}
	}
	return types.PaymentsSummaryResponse{
		Default: types.PaymentsSummary{
			TotalRequests: defaultCount,
			TotalAmount:   types.RoundedFloat(defaultSum),
		},
		Fallback: types.PaymentsSummary{
			TotalRequests: fallbackCount,
			TotalAmount:   types.RoundedFloat(fallbackSum),
		},
	}
}

func parseTimeRange(c adapter.RequestContext) (from, to time.Time, filter bool, err error) {
	getQuery := func(name string) string {
		return c.Query(name)
	}
	fromParam, toParam := getQuery("from"), getQuery("to")
	isEmpty := func(s string) bool { return stringx.IsEmpty(s) }
	if isEmpty(fromParam) && isEmpty(toParam) {
		return
	}
	if isEmpty(fromParam) || isEmpty(toParam) {
		err = http.ErrMissingFile
		return
	}
	from, err = date.Parse(fromParam)
	if err != nil {
		return
	}
	to, err = date.Parse(toParam)
	if err != nil {
		return
	}
	filter = true
	return
}
