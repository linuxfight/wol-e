package main

import (
	"flag"
	"github.com/go-co-op/gocron/v2"
	"os"
	"os/signal"
	"syscall"
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
	s, err := gocron.NewScheduler(gocron.WithLocation(time.Local))
	if err != nil {
		logger.Log.Panic(err)
	}

	for _, d := range cfg.Devices {
		err := d.InitCron(s)
		if err != nil {
			logger.Log.Panic(err)
		}
	}

	bot, err := telebot.NewBot(*cfg.BotConfig)
	if err != nil {
		logger.Log.Panic(err)
	}

	bot.Use(middleware.Whitelist(cfg.AdminIds...))
	bot.Use(middleware.Recover())
	// bot.Use(middleware.Logger())

	bot.Handle(telebot.OnText, func(context telebot.Context) error {
		return handlers.Message(context, cfg)
	})
	bot.Handle(telebot.OnCallback, func(context telebot.Context) error {
		return handlers.Callback(context, cfg)
	})
	logger.Log.Infof("bot started - https://t.me/%s", bot.Me.Username)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Log.Info("Starting bot...")
		bot.Start()
	}()

	// Block until we receive shutdown signal
	<-sigChan
	shutdown(s, bot)
}

func shutdown(s gocron.Scheduler, bot *telebot.Bot) {
	logger.Log.Info("Shutting down...")

	// Stop the bot first
	logger.Log.Info("Stopping bot...")
	bot.Stop()
	logger.Log.Info("Bot stopped successfully")

	// Stop scheduled jobs
	logger.Log.Info("Stopping scheduler jobs...")
	if err := s.StopJobs(); err != nil {
		logger.Log.Errorf("Error stopping jobs: %v", err)
	} else {
		logger.Log.Info("All jobs stopped")
	}

	// Shutdown scheduler
	logger.Log.Info("Shutting down scheduler...")
	if err := s.Shutdown(); err != nil {
		logger.Log.Errorf("Error shutting down scheduler: %v", err)
	} else {
		logger.Log.Info("Scheduler stopped successfully")
	}

	logger.Log.Info("Shutdown complete")
	os.Exit(0)
}
