package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gabrielluizsf/rinha-backend-2005/db"
	"github.com/gabrielluizsf/rinha-backend-2005/env"
	"github.com/gabrielluizsf/rinha-backend-2005/requests"
	"github.com/gabrielluizsf/rinha-backend-2005/types"
	"github.com/i9si-sistemas/nine"
	"github.com/i9si-sistemas/nine/pkg/client"
	"github.com/i9si-sistemas/stringx"
)

func processPayment(ctx context.Context, payment types.PaymentRequest) error {
	health, err := RetrieveHealthStates(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve health states: %w", err)
	}

	processorURL := env.Get().ProcessorDefaultURL
	processorType := "DEFAULT"
	status := "PROCESSED_DEFAULT"
	if health.DefaultProcessor.Failing || health.DefaultProcessor.MinResponseTime > health.FallBackProcessor.MinResponseTime+50 {
		processorURL = env.Get().ProcessorFallbackURL
		processorType = "FALLBACK"
		status = "PROCESSED_FALLBACK"
	}

	requestedAt := time.Now().UTC().Format(time.RFC3339Nano)

	body := nine.JSON{
		"correlationId": payment.CorrelationID,
		"amount":        payment.Amount,
		"requestedAt":   requestedAt,
	}

	jsonBody, _ := body.Buffer()

	res, err := requests.NewWithContext(ctx).Post(processorURL+"/payments", &client.Options{
		Headers: []client.Header{{Data: client.Data{Key: "Content-Type", Value: "application/json"}}},
		Body:    jsonBody,
	})
	if err != nil {
		return err
	} else {
		defer res.Body.Close()
		if res.StatusCode < 200 || res.StatusCode >= 300 {
			key := stringx.New("health:").Concat(stringx.New(processorType).ToLowerCase())
			db.HSet(key.String(), "failing", "true")
			return fmt.Errorf("failed to process payment: %s", res.Status)
		}
	}

	processedPayment := types.ProcessedPayment{
		CorrelationID: payment.CorrelationID,
		Amount:        payment.Amount,
		Status:        status,
		Processor:     processorType,
		CreatedAt:     requestedAt,
	}

	paymentData, err := json.Marshal(processedPayment)
	if err != nil {
		return fmt.Errorf("failed to marshal payment data: %w", err)
	}

	err = db.HSet("payments", payment.CorrelationID, paymentData)
	if err != nil {
		return fmt.Errorf("failed to save payment in redis: %w", err)
	}

	return nil
}
