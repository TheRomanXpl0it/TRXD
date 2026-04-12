package consts

import (
	"os"
	"strconv"
	"strings"

	"trxd/utils/log"
)

var AntiPanic = true

type Config struct {
	Name        string
	Value       any
	Type        string
	Category    string
	Description string
	Secret      bool
}

var DefaultConfigs map[string]Config = map[string]Config{
	"allow-register": {
		Name:        "Allow Register",
		Value:       false,
		Type:        "bool",
		Category:    "",
		Description: "whether to allow user registration",
		Secret:      false,
	},
	"chall-min-points": {
		Name:        "Challenge Min Points",
		Value:       50,
		Type:        "int",
		Category:    "",
		Description: "the minimum points a challenge can award",
		Secret:      false,
	},
	"chall-points-decay": {
		Name:        "Challenge Points Decay",
		Value:       15,
		Type:        "int",
		Category:    "",
		Description: "the rate at which challenge points decay",
		Secret:      false,
	},
	"instance-lifetime": {
		Name:        "Instance Lifetime",
		Value:       30 * 60, // 30 minutes
		Type:        "duration",
		Category:    "instances",
		Description: "the instances lifetime duration in seconds",
		Secret:      false,
	},
	"reclaim-instance-interval": {
		Name:        "Reclaim Instance Interval",
		Value:       5 * 60, // 5 minutes
		Type:        "duration",
		Category:    "instances",
		Description: "the interval for reclaiming instances in seconds",
		Secret:      false,
	},
	"instance-max-memory": {
		Name:        "Instance Max Memory",
		Value:       512,
		Type:        "int",
		Category:    "instances",
		Description: "the maximum memory allocation for each instance in MB",
		Secret:      false,
	},
	"instance-max-cpu": {
		Name:        "Instance Max CPU",
		Value:       1.0,
		Type:        "float",
		Category:    "instances",
		Description: "the maximum CPU allocation for each instance",
		Secret:      false,
	},
	"min-port": {
		Name:        "Min Port",
		Value:       10000,
		Type:        "port",
		Category:    "instances",
		Description: "the minimum port number for instance allocation",
		Secret:      false,
	},
	"max-port": {
		Name:        "Max Port",
		Value:       20000,
		Type:        "port",
		Category:    "instances",
		Description: "the maximum port number for instance allocation",
		Secret:      false,
	},
	"hash-len": {
		Name:        "Hash Domain Length",
		Value:       12,
		Type:        "int",
		Category:    "instances",
		Description: "the length of the random hash used into the instance domain (e.g. abcdef123456.domain.com)",
		Secret:      false,
	},
	"domain": {
		Name:        "Domain",
		Value:       "",
		Type:        "url",
		Category:    "instances",
		Description: "the domain used for the instances (e.g. domain.com)",
		Secret:      false,
	},
	"discord-webhook": {
		Name:        "Discord Webhook",
		Value:       "",
		Type:        "url",
		Category:    "",
		Description: "the Discord webhook URL for first blood notifications",
		Secret:      false,
	},
	"user-mode": {
		Name:        "Single User Mode",
		Value:       false,
		Type:        "bool",
		Category:    "",
		Description: "if enabled there will be no teams, but only users, like a single player mode",
		Secret:      false,
	},
	"scoreboard-top": {
		Name:        "Scoreboard Top Teams",
		Value:       10,
		Type:        "int",
		Category:    "",
		Description: "the number of the top teams to show on the scoreboard graph",
		Secret:      false,
	},
	"start-time": {
		Name:        "Start Time",
		Value:       "",
		Type:        "date",
		Category:    "",
		Description: "the start time for the competition following the RFC3339 format (e.g. 2024-01-02T15:04:05Z07:00)",
		Secret:      false,
	},
	"end-time": {
		Name:        "End Time",
		Value:       "",
		Type:        "date",
		Category:    "",
		Description: "the end time for the competition following the RFC3339 format (e.g. 2024-01-02T15:04:05Z07:00)",
		Secret:      false,
	},
	"email-verification": {
		Name:        "Email Verification",
		Value:       false,
		Type:        "bool",
		Category:    "email",
		Description: "enables all the email related features",
		Secret:      false,
	},
	"jwt-secret": {
		Name:        "Email JWT Secret",
		Value:       "",
		Type:        "string",
		Category:    "email",
		Description: "the secret key used for signing JWT tokens for email verification",
		Secret:      true,
	},
	"email-server": {
		Name:        "Email Server",
		Value:       "",
		Type:        "url",
		Category:    "email",
		Description: "the SMTP server address for sending verification emails (e.g. smtp.example.com)",
		Secret:      false,
	},
	"email-port": {
		Name:        "Email Server Port",
		Value:       587,
		Type:        "port",
		Category:    "email",
		Description: "the port number for the SMTP server (e.g. 587)",
		Secret:      false,
	},
	"email-addr": {
		Name:        "Email Address",
		Value:       "",
		Type:        "string",
		Category:    "email",
		Description: "the email address for sending verification emails (e.g. no-reply@example.com)",
		Secret:      false,
	},
	"email-passwd": {
		Name:        "Email Password",
		Value:       "",
		Type:        "string",
		Category:    "email",
		Description: "the password for the email account used for sending verification emails",
		Secret:      true,
	},
}

// var DefaultConfigs = map[string]any{
// 	"allow-register":            false,
// 	"chall-min-points":          50,
// 	"chall-points-decay":        15,
// 	"instance-lifetime":         30 * 60, // 30 minutes
// 	"reclaim-instance-interval": 5 * 60,  // 5 minutes
// 	"instance-max-memory":       512,
// 	"instance-max-cpu":          1.0,
// 	"min-port":                  10000,
// 	"max-port":                  20000,
// 	"hash-len":                  12,
// 	"domain":                    "",
// 	"discord-webhook":           "",
// 	"user-mode":                 false,
// 	"scoreboard-top":            10,
// 	"start-time":                "",
// 	"end-time":                  "",
// 	"jwt-secret":                "",
// 	"email-verification":        false,
// 	"email-server":              "",
// 	"email-port":                587,
// 	"email-addr":                "",
// 	"email-passwd":              "",
// }

func LoadEnvConfigs() {
	if os.Getenv("DISABLE_ANTI_PANIC") != "" {
		AntiPanic = false
		log.Warn("Anti Panic disabled")
	}

	for name, conf := range DefaultConfigs {
		envName := strings.ReplaceAll(name, "-", "_")
		envName = strings.ToUpper(envName)

		newValue := os.Getenv(envName)
		if newValue == "" {
			continue
		}

		switch conf.Value.(type) {
		case bool:
			if newValue == "1" || strings.ToLower(newValue) == "true" {
				conf.Value = true
			}
		case int:
			intValue, err := strconv.Atoi(newValue)
			if err != nil {
				log.Warn("Invalid int value for env", "env", envName, "value", newValue)
				continue
			}
			conf.Value = intValue
		case float64:
			floatValue, err := strconv.ParseFloat(newValue, 64)
			if err != nil {
				log.Warn("Invalid float value for env", "env", envName, "value", newValue)
				continue
			}
			conf.Value = floatValue
		case string:
			conf.Value = newValue
		default:
			log.Fatal("Unsupported config type for env", "env", envName, "type", conf.Type)
			continue
		}

		DefaultConfigs[name] = conf
	}
}
