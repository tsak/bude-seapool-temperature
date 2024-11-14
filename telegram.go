package main

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log/slog"
	"net/http"
	"time"
)

func StartTelegramBot(sm *StateManager, monnit *Monnit, cfg *Config) error {
	options := []bot.Option{
		bot.WithHTTPClient(5*time.Second, &http.Client{
			Timeout: 5 * time.Second,
		}),
	}
	if cfg.Debug {
		options = append(options, bot.WithDebug())
	}

	b, err := bot.New(cfg.TelegramToken)
	if err != nil {
		return err
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, startHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, startHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/temp", bot.MatchTypeExact, getTemperatureHandler(sm, monnit))

	b.Start(context.TODO())

	return nil
}

func startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "I understand the /temp command",
	})
}

// TelegramBot listens to messages on an updates channel and replies to bot commands
func getTemperatureHandler(sm *StateManager, monnit *Monnit) func(ctx context.Context, b *bot.Bot, update *models.Update) {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		last := monnit.LastReading()
		text := fmt.Sprintf("Temperature was <b>%s</b> at %s", last.Temperature.String(), last.MessageDate.String())
		msg := bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      text,
			ParseMode: models.ParseModeHTML,
		}
		if _, err := b.SendMessage(ctx, &msg); err != nil {
			slog.Warn("unable to respond to telegram command", "command", "/temp", "error", err)
		}
		sm.IncrementBotRequests()
	}
}
