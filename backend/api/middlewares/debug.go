package middlewares

import (
	"encoding/json"

	"trxd/utils/log"

	"github.com/gofiber/fiber/v2"
)

func decodeJSONBody(body []byte) string {
	if len(body) == 0 {
		return ""
	}

	defaultBody := string(body)

	var tmp interface{}
	if err := json.Unmarshal(body, &tmp); err == nil {
		if tmp2, err := json.MarshalIndent(tmp, "", "  "); err == nil {
			return string(tmp2)
		}
	}

	return defaultBody
}

func Debug(c *fiber.Ctx) error {
	if c.Path() == "/countries" {
		return c.Next()
	}

	reqBody := c.BodyRaw()
	body := decodeJSONBody(reqBody)
	log.Debug("Request:", "method", c.Method(), "path", c.Path(), "body", body)

	e := c.Next()

	resStatus := c.Response().StatusCode()
	resBody := c.Response().Body()
	body = decodeJSONBody(resBody)
	log.Debug("Response:", "status", resStatus, "body", body)

	return e
}
