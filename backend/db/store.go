package db

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	fiberRedis "github.com/gofiber/storage/redis/v3"
	"github.com/redis/go-redis/v9"
	"github.com/tde-nico/log"
)

var rdb *redis.Client
var storage map[string]string
var storageRWMutex sync.RWMutex
var Store *session.Store

func initStorage(host string, port int, password string) {
	storeConf := session.Config{
		Expiration:     30 * 24 * time.Hour,
		CookiePath:     "/",
		CookieSameSite: fiber.CookieSameSiteLaxMode,
	}

	storage = make(map[string]string)
	if os.Getenv("REDIS_DISABLE") == "" {
		rdb = redis.NewClient(&redis.Options{
			Addr:      fmt.Sprintf("%s:%d", host, port),
			Password:  password,
			DB:        0,
			TLSConfig: nil,
		})
		storeConf.Storage = fiberRedis.NewFromConnection(rdb)
	}

	Store = session.New(storeConf)
}

func StorageSet(ctx context.Context, key string, val string) error {
	if rdb == nil {
		storageRWMutex.Lock()
		defer storageRWMutex.Unlock()
		storage[key] = val
		return nil
	}

	err := rdb.Set(ctx, key, []byte(val), 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func StorageSetNX(ctx context.Context, key string, val string, expiration ...time.Duration) (bool, error) {
	if rdb == nil {
		log.Warn("Expiration parameter is ignored in in-memory storage")
		storageRWMutex.Lock()
		defer storageRWMutex.Unlock()
		if _, ok := storage[key]; ok {
			return false, nil
		}
		storage[key] = val
		return true, nil
	}

	exp := 0 * time.Second
	if len(expiration) > 0 {
		exp = expiration[0]
	}

	res, err := rdb.SetArgs(ctx, key, []byte(val), redis.SetArgs{
		Mode: "NX",
		TTL:  exp,
	}).Result()
	if err != nil {
		return false, err
	}

	return res == "OK", nil
}

func StorageGet(ctx context.Context, key string) (*string, error) {
	if rdb == nil {
		storageRWMutex.RLock()
		defer storageRWMutex.RUnlock()
		val, ok := storage[key]
		if !ok {
			return nil, nil
		}
		return &val, nil
	}

	val, err := rdb.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	if val == nil {
		return nil, nil
	}

	strVal := string(val)
	return &strVal, nil
}

func StorageDelete(ctx context.Context, key string) error {
	if rdb == nil {
		storageRWMutex.Lock()
		defer storageRWMutex.Unlock()
		delete(storage, key)
		return nil
	}

	err := rdb.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}

func StorageFlush(ctx context.Context) error {
	if rdb == nil {
		storageRWMutex.Lock()
		defer storageRWMutex.Unlock()
		storage = make(map[string]string)
		return nil
	}

	err := rdb.FlushAll(ctx).Err()
	if err != nil {
		return err
	}

	return nil
}
