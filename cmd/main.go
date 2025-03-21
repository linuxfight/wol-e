package main

import (
	"flag"
	"time"
	"wol-e/internal/config"
	"wol-e/internal/handlers"
	"wol-e/internal/logger"

	"gopkg.in/telebot.v4"
	"gopkg.in/telebot.v4/middleware"
)

func main() {
	localTime := time.Now()
	timezone, _ := localTime.Zone()

	configPath := flag.String("config", "./config.yaml", "path to a config file, ex. /etc/wol-e/config.yaml")
	flag.Parse()

	logger.New(false, timezone)
	cfg := config.New(*configPath)

	bot, err := telebot.NewBot(*cfg.BotConfig)
	if err != nil {
		logger.Log.Panicf("error creating new bot: %v", err)
	}

	bot.Use(middleware.Whitelist())

	bot.Handle(telebot.OnText, func(context telebot.Context) error {
		return handlers.Message(context, cfg)
	})

	bot.Handle(telebot.OnCallback, func(context telebot.Context) error {
		return handlers.Callback(context, cfg)
	})

	logger.Log.Infof("bot started - https://t.me/%s", bot.Me.Username)
	bot.Start()
}
