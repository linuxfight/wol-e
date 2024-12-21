package config

import (
	"github.com/spf13/viper"
	"gopkg.in/telebot.v4"
	"time"
	"wol-e/internal/device"
	"wol-e/internal/logger"
)

type AppConfig struct {
	BotConfig *telebot.Settings
	Devices   []device.Device
	AdminId   int64
}

var (
	Config *AppConfig
)

func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		logger.Log.Panicf("failed to read config: %v", err)
		return
	}

	var devices []device.Device
	if err := viper.UnmarshalKey("devices", &devices); err != nil {
		logger.Log.Panicf("error unmarshaling data: %v", err)
		return
	}

	Config = &AppConfig{
		BotConfig: &telebot.Settings{
			Token:   viper.GetString("bot.token"),
			Verbose: viper.GetBool("settings.debug"),
			Poller: &telebot.LongPoller{
				Timeout:        10 * time.Second,
				AllowedUpdates: []string{"message", "callback_query"},
			},
		},
		Devices: devices,
		AdminId: viper.GetInt64("bot.admin_id"),
	}
}
