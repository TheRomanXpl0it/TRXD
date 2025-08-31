package discord

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"trxd/db"
	"trxd/db/sqlc"

	"github.com/tde-nico/log"
)

const WebhookTimeout = 5 * time.Second

func BroadcastWebhook(url string, body interface{}) error {
	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	client := &http.Client{Timeout: WebhookTimeout}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("discord returned status: %s", resp.Status)
	}

	return nil
}

func BroadcastFirstBlood(ctx context.Context, challenge *sqlc.Challenge, uid int32) {
	conf, err := db.GetConfig(ctx, "discord-webhook")
	if err != nil {
		log.Error("Failed to fetch webhook url:", "err", err)
		return
	}

	if conf == nil || conf.Value == "" {
		return
	}

	team, err := db.GetTeamFromUser(ctx, uid)
	if err != nil {
		log.Error("Failed to fetch user's team:", "err", err)
		return
	}

	// TODO: Hardcoded format:(
	name := strings.ReplaceAll(team.Name, "@", "@\u200b")
	msg := fmt.Sprintf("First blood for **%s** goes to **%s**! 🩸", challenge.Name, name)
	body := map[string]string{"content": msg}

	if err := BroadcastWebhook(conf.Value, body); err != nil {
		log.Error("Failed to send webhook:", "err", err)
	}
}
