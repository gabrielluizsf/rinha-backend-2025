package db

import (
	"time"
)

func Set(
	key string, value any,
	expiration time.Duration,
) error {
	return client.Set(ctx, key, value, expiration).Err()
}

type GetAllResult map[string]string

func Save(key string, data any) error {
	return client.LPush(ctx, key, data).Err()
}

func Purge(key string) error {
	return client.Del(ctx, key).Err()
}

func GetAll(key string) (GetAllResult, error) {
	return client.HGetAll(ctx, key).Result()
}

func Get(key string) (value string, err error) {
	return client.Get(ctx, key).Result()
}

func HSet(key string, values ...any) error {
	return client.HSet(ctx, key, values...).Err()
}

func RPopLPush(
	source,
	dest string,
) (string, error) {
	return client.RPopLPush(ctx, source, dest).Result()
}

func SetNX(
	key, value string,
	expiration time.Duration,
) (ok bool, err error) {
	return client.SetNX(ctx, key, value, expiration).Result()
}
