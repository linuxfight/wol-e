package main

import (
	"flag"
	"gopkg.in/telebot.v4"
	"time"
	"wol-e/internal/config"
	"wol-e/internal/handlers"
	"wol-e/internal/logger"
)

func main() {
	localTime := time.Now()
	timezone, _ := localTime.Zone()

	configPath := flag.String("config", "./config.yaml", "path to a config file, ex. /etc/wol-e/config.yaml")
	flag.Parse()

	logger.New(false, timezone)
	config.New(*configPath)

	bot, err := telebot.NewBot(*config.Config.BotConfig)
	if err != nil {
		logger.Log.Panicf("error creating new bot: %v", err)
	}

	bot.Handle(telebot.OnText, func(context telebot.Context) error {
		return handlers.Message(context)
	})

	bot.Handle(telebot.OnCallback, func(context telebot.Context) error {
		return handlers.Callback(context)
	})

	logger.Log.Infof("bot started - https://t.me/%s", bot.Me.Username)
	bot.Start()
}
