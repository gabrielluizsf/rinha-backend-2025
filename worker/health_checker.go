package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gabrielluizsf/rinha-backend-2005/db"
	"github.com/gabrielluizsf/rinha-backend-2005/env"
	"github.com/gabrielluizsf/rinha-backend-2005/requests"
	"github.com/gabrielluizsf/rinha-backend-2005/types"
	"github.com/i9si-sistemas/nine/pkg/client"
)

var isLeader atomic.Bool

func StartLeaderElection() {
	instanceID := env.Get().InstanceID
	if instanceID == "" {
		instanceID = fmt.Sprintf("instance-%d", time.Now().UnixNano())
	}

	go func() {
		for {
			ok, err := db.SetNX("rinha-leader-lock", instanceID, 10*time.Second)
			if err != nil {
				fmt.Println("Redis error during leader election:", err)
			}
			if ok {
				fmt.Println("Became leader:", instanceID)
				isLeader.Store(true)
				go renewLeaderLock(instanceID)
				go healthChecker()
				return
			} else {
				isLeader.Store(false)
			}

			time.Sleep(3 * time.Second)
		}
	}()
}

func renewLeaderLock(instanceID string) {
	ticker := time.NewTicker(8 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		val, err := db.Get("rinha-leader-lock")
		if err == nil && val == instanceID {
			db.Set("rinha-leader-lock", instanceID, 10*time.Second)
		} else {
			fmt.Println("Lost leadership")
			isLeader.Store(false)
			return
		}
	}
}

func healthChecker() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	processorURLs := []string{
		env.Get().ProcessorDefaultURL,
		env.Get().ProcessorFallbackURL,
	}

	redisKeys := []string{"health:default", "health:fallback"}

	for {
		<-ticker.C
		for i, url := range processorURLs {
			go func(i int, url string) {
				res, err := requests.NewWithContext(context.TODO()).Get(url + "/payments/service-health", &client.Options{})
				if err != nil {
					fmt.Printf("Health check failed for %s: %v\n", url, err)
					return
				}
				defer res.Body.Close()

				if res.StatusCode != http.StatusOK {
					fmt.Printf("Health check non-200 for %s: %d\n", url, res.StatusCode)
					return
				}

				var hr types.HealthResponse

				if err := json.NewDecoder(res.Body).Decode(&hr); err != nil {
					fmt.Printf("Failed to decode health response from %s: %v\n", url, err)
					return
				}

				fmt.Printf("Health for %s: failing=%v, minResponseTime=%d\n", url, hr.Failing, hr.MinResponseTime)

				data, err := json.Marshal(hr)
				if err != nil {
					fmt.Printf("Failed to marshal health response for %s: %v\n", url, err)
					return
				}

				if err := db.Set(redisKeys[i], data, 0); err != nil {
					fmt.Printf("Failed to save health state for %s in Redis: %v\n", url, err)
				}
			}(i, url)
		}
	}
}

func RetrieveHealthStates(ctx context.Context) (*types.HealthManager, error) {
	defaultKey := "health:default"
	fallbackKey := "health:fallback"

	defaultVal, err := db.Get(defaultKey)
	if err != nil {
		return nil, err
	}
	fallbackVal, err := db.Get(fallbackKey)
	if err != nil {
		return nil, err
	}

	var defaultHealth types.HealthResponse
	var fallbackHealth types.HealthResponse

	if err := json.Unmarshal([]byte(defaultVal), &defaultHealth); err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(fallbackVal), &fallbackHealth); err != nil {
		return nil, err
	}

	return &types.HealthManager{
		DefaultProcessor:  defaultHealth,
		FallBackProcessor: fallbackHealth,
	}, nil
}
