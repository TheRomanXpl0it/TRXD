package db

import (
	"os"
	"time"
	"trxd/utils/consts"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis/v3"
	"github.com/tde-nico/log"
)

var RedisStorage *redis.Storage
var Store *session.Store

func initStorage(host string, port int, password string) {
	storeConf := session.Config{
		Expiration:     30 * 24 * time.Hour,
		CookiePath:     "/",
		CookieSameSite: fiber.CookieSameSiteLaxMode,
	}

	if os.Getenv("REDIS_DISABLED") == "" {
		RedisStorage = redis.New(redis.Config{
			Host:      host,
			Port:      port,
			Password:  password,
			Database:  0,
			Reset:     false,
			TLSConfig: nil,
		})
		storeConf.Storage = RedisStorage
	} else if !consts.Testing {
		log.Warn("Redis storage disabled")
	}

	Store = session.New(storeConf)
}
