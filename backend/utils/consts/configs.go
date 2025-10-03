package consts

var Testing = false

var DefaultConfigs = map[string]any{
	"allow-register":            false,
	"chall-min-points":          50,
	"chall-points-decay":        15,
	"instance-lifetime":         30 * 60, // 30 minutes
	"reclaim-instance-interval": 5 * 60,  // 5 minutes
	"instance-max-memory":       512,
	"instance-max-cpu":          "1.0",
	"min-port":                  20000,
	"max-port":                  30000,
	"hash-len":                  12,
	"secret":                    "",
	"domain":                    "",
	"discord-webhook":           "",
	"project-name":              "trxd",
}
