package db

import (
	"context"
	"fmt"
	"os"
	"time"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	fiberRedis "github.com/gofiber/storage/redis/v3"
	"github.com/redis/go-redis/v9"
	"github.com/tde-nico/log"
)

var rdb *redis.Client
var Store *session.Store

func initStorage(host string, port int, password string) {
	storeConf := session.Config{
		Expiration:     30 * 24 * time.Hour,
		CookiePath:     "/",
		CookieSameSite: fiber.CookieSameSiteLaxMode,
	}

	if os.Getenv("REDIS_DISABLE") == "" {
		rdb = redis.NewClient(&redis.Options{
			Addr:      fmt.Sprintf("%s:%d", host, port),
			Password:  password,
			DB:        0,
			TLSConfig: nil,
		})
		storeConf.Storage = fiberRedis.NewFromConnection(rdb)
	} else if !consts.Testing {
		log.Warn("Redis storage disabled")
	}

	Store = session.New(storeConf)
}

func StorageSet(ctx context.Context, key string, val string) error {
	if rdb == nil {
		return nil
	}

	err := rdb.Set(ctx, key, []byte(val), 0).Err()
	if err != nil {
		return err
	}

	return nil
}

// var mapSetNX map[string]string
// var lockMapSetNX sync.RWMutex

// func StorageSetNX(ctx context.Context, key string, val string) (bool, error) {
// 	if rdb == nil {
// 		lockMapSetNX.Lock()
// 		defer lockMapSetNX.Unlock()
// 		if _, ok := mapSetNX[key]; ok {
// 			return false, nil
// 		}
// 		mapSetNX[key] = val
// 		return true, nil
// 	}

// 	res, err := rdb.SetNX(ctx, key, []byte(val), 0).Result()
// 	if err != nil {
// 		return false, err
// 	}

// 	return res, nil
// }

func StorageGet(ctx context.Context, key string) (*string, error) {
	if rdb == nil {
		return nil, nil
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

func StorageFlush() error {
	if rdb == nil {
		return nil
	}

	err := rdb.FlushAll(context.Background()).Err()
	if err != nil {
		return err
	}

	return nil
}
