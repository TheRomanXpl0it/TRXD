package discord

import (
	"fmt"
	"time"
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"trxd/db"
	"trxd/db/sqlc"

	"github.com/tde-nico/log"
)

const WebhookTimeout = 5 * time.Second

type WebhookBody struct {
	Content string `json:"content"`
}

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
		log.Error("Failed to fetch webhook url: %v", err)
		return
	}

	// Discord broadcast is disabled
	if conf == nil {
		return
	}

	team, err := db.GetTeamFromUser(ctx, uid)
	if err != nil {
		log.Error("Failed to fetch user's team: %v", err)
		return
	}

	// TODO: it's hardcoded :(
	msg := fmt.Sprintf("First blood for **%s** goes to **%s**! 🩸", challenge.Name, team.Name)
	body := WebhookBody{
		Content: msg,
	}

	if err := BroadcastWebhook(conf.Value, body); err != nil {
		log.Error("Failed to send webhook: %v", err)
	}
}
