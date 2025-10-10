package db

import (
	"context"
	"os"
	"time"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis/v3"
	"github.com/tde-nico/log"
)

var redisStorage *redis.Storage
var Store *session.Store

func initStorage(host string, port int, password string) {
	storeConf := session.Config{
		Expiration:     30 * 24 * time.Hour,
		CookiePath:     "/",
		CookieSameSite: fiber.CookieSameSiteLaxMode,
	}

	if os.Getenv("REDIS_DISABLE") == "" {
		redisStorage = redis.New(redis.Config{
			Host:      host,
			Port:      port,
			Password:  password,
			Database:  0,
			Reset:     false,
			TLSConfig: nil,
		})
		storeConf.Storage = redisStorage
	} else if !consts.Testing {
		log.Warn("Redis storage disabled")
	}

	Store = session.New(storeConf)
}

func StorageSet(ctx context.Context, key string, val string) error {
	if redisStorage == nil {
		return nil
	}

	if val == "" {
		val = "\x00"
	}
	err := redisStorage.SetWithContext(ctx, key, []byte(val), 0)
	if err != nil {
		return err
	}

	return nil
}

func StorageGet(ctx context.Context, key string) (*string, error) {
	if redisStorage == nil {
		return nil, nil
	}

	val, err := redisStorage.GetWithContext(ctx, key)
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, nil
	}

	strVal := string(val)
	if strVal == "\x00" {
		strVal = ""
	}

	return &strVal, nil
}

func StorageFlush() error {
	if redisStorage == nil {
		return nil
	}

	err := redisStorage.Reset()
	if err != nil {
		return err
	}

	return nil
}
