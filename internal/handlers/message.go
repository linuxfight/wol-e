package handlers

import (
	"fmt"
	"gopkg.in/telebot.v4"
	"wol-e/internal/config"
	"wol-e/internal/logger"
)

func handleStart(context telebot.Context) error {
	return context.Send("Hello, I'm your WoL bot! Type /help to get started!")
}

func handleHelp(context telebot.Context) error {
	text := "Here's the list of commands:" + "\n" +
		"/start - start the bot" + "\n" +
		"/help - show this list" + "\n" +
		"/devices - show all devices"
	return context.Send(text)
}

func handleDevices(context telebot.Context) error {
	for i, device := range config.Config.Devices {
		text := device.GenerateBotText()
		replyMarkup := &telebot.ReplyMarkup{}
		pingButton := telebot.Btn{
			Text: "🔄",
			Data: fmt.Sprintf("ping:%d", i),
		}
		powerButton := telebot.Btn{
			Text: "🔌",
			Data: fmt.Sprintf("power:%d", i),
		}
		replyMarkup.Inline(replyMarkup.Row(pingButton, powerButton))
		//replyMarkup.Inline(replyMarkup.Row(powerButton))
		if err := context.Send(text, replyMarkup); err != nil {
			return err
		}
	}
	return nil
}

func Message(context telebot.Context) error {
	message := context.Message()

	if message.Chat.ID != config.Config.AdminId ||
		message.Sender.ID != config.Config.AdminId {
		logger.Log.Errorf("unauthorized user - %d - %d", message.Chat.ID, message.Sender.ID)
		return nil
	}

	switch context.Text() {
	case "/start":
		return handleStart(context)
	case "/help":
		return handleHelp(context)
	case "/devices":
		return handleDevices(context)
	}

	return nil
}
