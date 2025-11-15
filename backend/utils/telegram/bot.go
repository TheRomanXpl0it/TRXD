package telegram

import (
	"fmt"
	"context"
	"strconv"

	"trxd/db"
	"trxd/db/sqlc"

	"github.com/tde-nico/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func BroadcastMessage(bot *tgbotapi.BotAPI, chatId int64, message string) error {
	msg := tgbotapi.NewMessage(chatId, message)
	_, err := bot.Send(msg)
	return err
}

func BroadcastAnimation(bot *tgbotapi.BotAPI, chatId int64, animation string) error {
	msg := tgbotapi.NewAnimationShare(chatId, animation)
	_, err := bot.Send(msg)
	return err
}

func BroadcastSticker(bot *tgbotapi.BotAPI, chatId int64, sticker string) error {
	msg := tgbotapi.NewStickerShare(chatId, sticker)
	_, err := bot.Send(msg)
	return err
}

func BroadcastFirstBlood(ctx context.Context, challenge *sqlc.Challenge, uid int32) {
	token, err := db.GetConfig(ctx, "telegram-token")
	if err != nil {
		log.Error("Failed to fetch bot token:", "err", err)
		return
	}
	if token == "" {
		return
	}

	conf, err := db.GetConfig(ctx, "telegram-chat-id")
	if err != nil {
		log.Error("Failed to fetch chat id:", "err", err)
		return
	}

	chatId, err := strconv.ParseInt(conf, 10, 64)
	if err != nil || chatId < 0 {
		return
	}

	team, err := db.GetTeamFromUser(ctx, uid)
	if err != nil {
		log.Error("Failed to fetch user's team:", "err", err)
		return
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Error("Failed to create bot:", "err", err)
		return
	}

	// TODO: not hardcoded format (maybe into a config)
	challengeName := challenge.Name
	teamName := team.Name
	msg := fmt.Sprintf("First blood for **%s** goes to **%s**! 🩸", challengeName, teamName)

	if err := BroadcastMessage(bot, chatId, msg); err != nil {
		log.Error("Failed to send message:", "err", err)
		return
	}

	sticker := "CAACAgQAAxkBAAExXY1nn3i7L9e7bC_Pt5vHj6wpTfL4WwACyRcAAiog-VAHs5wSbsrARDYE"

	if err := BroadcastSticker(bot, chatId, sticker); err != nil {
		log.Error("Failed to send media:", "err", err)
	}
}
