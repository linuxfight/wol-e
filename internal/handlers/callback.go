package handlers

import (
	"gopkg.in/telebot.v4"
	"strconv"
	"strings"
	"wol-e/internal/config"
	"wol-e/internal/device"
	"wol-e/internal/logger"
)

func handlePower(context telebot.Context, target device.Device, callback *telebot.Callback) error {
	if err := target.TurnOn(); err != nil {
		return err
	}
	return context.Respond(&telebot.CallbackResponse{
		CallbackID: callback.ID,
		Text:       "turning on",
		ShowAlert:  false,
	})
}

/*
func handlePing(context telebot.Context, target device.Device, targetId int, callback *telebot.Callback) error {
	text := target.GenerateBotText()

	replyMarkup := &telebot.ReplyMarkup{}
	pingButton := telebot.Btn{
		Text: "🔄",
		Data: fmt.Sprintf("ping:%d", targetId),
	}
	powerButton := telebot.Btn{
		Text: "🔌",
		Data: fmt.Sprintf("power:%d", targetId),
	}
	replyMarkup.Inline(replyMarkup.Row(pingButton, powerButton))

	if err := context.Edit(text); err != nil {
		logger.Log.Errorf("failed to edit message: %v", err)
		return err
	}

	if err := context.Edit(replyMarkup); err != nil {
		logger.Log.Errorf("failed to edit message: %v", err)
		return err
	}

	return context.Respond(&telebot.CallbackResponse{
		CallbackID: callback.ID,
		Text:       "updated",
		ShowAlert:  false,
	})
}
*/

func Callback(context telebot.Context) error {
	callback := context.Callback()

	if callback.Message.Chat.ID != config.Config.AdminId ||
		callback.Sender.ID != config.Config.AdminId {
		logger.Log.Errorf("unauthorized user - %d - %d",
			callback.Message.Chat.ID, callback.Sender.ID)
		return nil
	}

	data := strings.Split(callback.Data, ":")
	method := data[0]
	deviceIdStr := data[1]
	deviceId, err := strconv.Atoi(deviceIdStr)
	if err != nil {
		logger.Log.Errorf("invalid user id %s", deviceIdStr)
		return nil
	}

	var target device.Device
	found := false
	for i, d := range config.Config.Devices {
		if i != deviceId {
			continue
		}
		target = d
		found = true
	}
	if !found {
		logger.Log.Errorf("invalid device id %s", deviceIdStr)
		return nil
	}

	switch method {
	/*
		case "ping":
			return handlePing(context, target, deviceId, callback)
	*/
	case "power":
		return handlePower(context, target, callback)
	}

	return nil
}
