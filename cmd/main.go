package main

import (
	"gopkg.in/telebot.v4"
	"time"
	"wol-e/internal/config"
	"wol-e/internal/logger"
)

func main() {
	localTime := time.Now()
	timezone, _ := localTime.Zone()

	logger.New(true, timezone)

	config.Init()

	bot, err := telebot.NewBot(*config.Config.BotConfig)
	if err != nil {
		logger.Log.Panicf("error creating new bot: %v", err)
	}

	bot.Handle(telebot.OnText, func(context telebot.Context) error {
		if context.Chat().ID != config.Config.AdminId || context.Sender().ID != config.Config.AdminId {
			logger.Log.Errorf("unauthorized user - %d - %d", context.Chat().ID, context.Sender().ID)
		}

		switch context.Text() {
		case "/start":
			return context.Send("Hello, I'm your WoL bot! Type /help to get started!")
		}

		return nil
	})

	logger.Log.Infof("bot started - https://t.me/%s", bot.Me.Username)
	bot.Start()
}
