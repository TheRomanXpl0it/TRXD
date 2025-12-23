package consts

import (
	"os"
	"strconv"
	"strings"

	"github.com/tde-nico/log"
)

var Testing = false

var DefaultConfigs = map[string]any{
	"allow-register":            false,
	"chall-min-points":          50,
	"chall-points-decay":        15,
	"instance-lifetime":         30 * 60, // 30 minutes
	"reclaim-instance-interval": 5 * 60,  // 5 minutes
	"instance-max-memory":       512,
	"instance-max-cpu":          1.0,
	"min-port":                  10000,
	"max-port":                  20000,
	"hash-len":                  12,
	"domain":                    "",
	"discord-webhook":           "",
	"project-name":              "trxd",
	"user-mode":                 false,
	"scoreboard-top":            10,
	"start-time":                "",
	"end-time":                  "",
}

func LoadEnvConfigs() {
	if os.Getenv("TESTING") != "" {
		Testing = true
		log.Warn("TESTING mode enabled")
	}

	for name, value := range DefaultConfigs {
		envName := strings.ReplaceAll(name, "-", "_")
		envName = strings.ToUpper(envName)

		newValue := os.Getenv(envName)
		if newValue == "" {
			continue
		}

		switch value.(type) {
		case bool:
			if newValue == "1" || strings.ToLower(newValue) == "true" {
				DefaultConfigs[name] = true
			}
		case int:
			intValue, err := strconv.Atoi(newValue)
			if err != nil {
				log.Warn("Invalid int value for env", "env", envName, "value", newValue)
				continue
			}
			DefaultConfigs[name] = intValue
		case float64:
			floatValue, err := strconv.ParseFloat(newValue, 64)
			if err != nil {
				log.Warn("Invalid float value for env", "env", envName, "value", newValue)
				continue
			}
			DefaultConfigs[name] = floatValue
		case string:
			DefaultConfigs[name] = newValue
		}

	}
}
