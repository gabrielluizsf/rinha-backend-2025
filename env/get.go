package env

import (
	"os"

	"github.com/i9si-sistemas/stringx"
)

type Env struct {
	Redis                string
	InstanceID           string
	ProcessorDefaultURL  string
	ProcessorFallbackURL string
	MaxWorkers           int
}

func Get() *Env {
	getEnv := func(key string) string { return os.Getenv(key) }
	workerConcurrency := getEnv("WORKER_CONCURRENCY")
	env := &Env{
		Redis:                getEnv("REDIS_URL"),
		InstanceID:           getEnv("INSTANCE_ID"),
		ProcessorDefaultURL:  getEnv("PROCESSOR_DEFAULT_URL"),
		ProcessorFallbackURL: getEnv("PROCESSOR_FALLBACK_URL"),
		MaxWorkers:           10,
	}
	if !stringx.IsEmpty(workerConcurrency) {
		maxWorkers, err := stringx.NewParser(workerConcurrency).Int()
		if err == nil && maxWorkers > 0 {
			env.MaxWorkers = int(maxWorkers)
		}
	}
	return env
}
