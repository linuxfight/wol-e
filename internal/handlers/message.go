package handlers

import (
	"fmt"
	"wol-e/internal/config"

	"gopkg.in/telebot.v4"
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

func handleDevices(context telebot.Context, cfg *config.Data) error {
	for i, device := range cfg.Devices {
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

func Message(context telebot.Context, cfg *config.Data) error {
	switch context.Text() {
	case "/start":
		return handleStart(context)
	case "/help":
		return handleHelp(context)
	case "/devices":
		return handleDevices(context, cfg)
	}

	return nil
}
