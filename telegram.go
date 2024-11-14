package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
)

// TelegramBot listens to messages on an updates channel and replies to bot commands
func TelegramBot(bot *tgbotapi.BotAPI, sm *StateManager, monnit *Monnit) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}

		response := tgbotapi.NewMessage(update.Message.Chat.ID, "Hello")

		last := monnit.LastReading()
		switch update.Message.Command() {
		case "temp":
			response.Text = fmt.Sprintf("Temperature was %s at %s", last.Temperature.String(), last.MessageDate.String())
		default:
			response.Text = "I understand /temp"
		}

		_, err := bot.Send(response)
		if err != nil {
			slog.Error("unable to send Telegram message", "error", err)
		}

		sm.IncrementBotRequests()
	}
}
