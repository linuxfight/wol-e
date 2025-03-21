package config

import (
	"fmt"
	"time"
	"wol-e/internal/device"
	"wol-e/internal/logger"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"gopkg.in/telebot.v4"
)

type Data struct {
	BotConfig *telebot.Settings
	Devices   []device.Device
	AdminIds  []int64
}

func New(configPath string) *Data {
	logger.Log.Infof("reading config file: %s", configPath)

	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		logger.Log.Panicf("failed to read config: %v", err)
		return nil
	}

	var devices []device.Device
	if err := viper.UnmarshalKey("devices", &devices); err != nil {
		logger.Log.Panicf("error unmarshaling data: %v", err)
		return nil
	}

	admins, err := getInt64Slice("bot.admins")
	if err != nil {
		logger.Log.Panicf("error unmarshaling data: %v", err)
		return nil
	}

	return &Data{
		BotConfig: &telebot.Settings{
			Token:   viper.GetString("bot.token"),
			Verbose: viper.GetBool("settings.debug"),
			Poller: &telebot.LongPoller{
				Timeout:        10 * time.Second,
				AllowedUpdates: []string{"message", "callback_query"},
			},
		},
		Devices:  devices,
		AdminIds: admins,
	}
}

// getInt64Slice retrieves a slice of int64 values from Viper for the specified key.
// Returns an error if the key does not exist or the value cannot be cast to []int64.
func getInt64Slice(key string) ([]int64, error) {
	val := viper.Get(key)
	if val == nil {
		return nil, fmt.Errorf("key '%s' not found", key)
	}

	// First convert to generic slice
	slice, err := cast.ToSliceE(val)
	if err != nil {
		return nil, fmt.Errorf("value at '%s' is not a slice: %v", key, err)
	}

	// Convert each element to int64
	result := make([]int64, 0, len(slice))
	for i, item := range slice {
		intVal, err := cast.ToInt64E(item)
		if err != nil {
			return nil, fmt.Errorf(
				"element %d at key '%s' cannot be converted to int64: %v",
				i, key, err,
			)
		}
		result = append(result, intVal)
	}

	return result, nil
}
