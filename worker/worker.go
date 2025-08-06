package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gabrielluizsf/rinha-backend-2005/db"
	"github.com/gabrielluizsf/rinha-backend-2005/env"
	"github.com/gabrielluizsf/rinha-backend-2005/types"
	"github.com/i9si-sistemas/stringx"
)

func StartWorker() {
	workerID := env.Get().InstanceID
	if stringx.IsEmpty(workerID) {
		workerID = fmt.Sprintf("worker-%d", time.Now().UnixNano())
	}
	processingQueue := "payments_processing:" + workerID

	for i := 0; i < env.Get().MaxWorkers; i++ {
		go func(workerNum int) {
			for {
				res, err := db.RPopLPush("payments_pending", processingQueue)
				if err != nil {
					time.Sleep(1 * time.Second)
					continue
				}

				var payment types.PaymentRequest

				if err := json.Unmarshal([]byte(res), &payment); err != nil {
					continue
				}

				if err := processPayment(context.Background(), payment); err != nil {
					db.Save("payments_pending", res)
				}
			}
		}(i)
	}
}
