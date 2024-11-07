package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"os"
)

func initTelegramBot() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file ", err)
		os.Exit(1)
	}
	botToken := os.Getenv("BOT_TOKEN")
	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	updates, _ := bot.UpdatesViaLongPolling(nil)
	defer bot.StopLongPolling()

	// Create a keyboard
	keyboard := tu.Keyboard(
		tu.KeyboardRow(
			tu.KeyboardButton("Get praying time by location").WithRequestLocation(),
			tu.KeyboardButton("Get praying time by city"),
		),
	).WithResizeKeyboard().WithInputFieldPlaceholder("Select something...")

	// Create a bot handler
	bh, _ := th.NewBotHandler(bot, updates)
	defer bh.Stop()

	// Handle /start command
	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		_, _ = bot.SendMessage(tu.Message(tu.ID(update.Message.Chat.ID), fmt.Sprintf("Hello %s! Welcome to this bot!", update.Message.From.FirstName)).WithReplyMarkup(keyboard))
	}, th.CommandEqual("start"))

	// Handle getting location
	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		_, _ = bot.SendMessage(tu.Message(tu.ID(update.Message.Chat.ID), "Send me your location").WithReplyMarkup(keyboard))
	}, th.TextEqual("Get praying time by location"))

	// Handle getting city
	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		_, _ = bot.SendMessage(tu.Message(tu.ID(update.Message.Chat.ID), "Send me your city").WithReplyMarkup(keyboard))
	}, th.TextEqual("Get praying time by city"))

	// Start handling updates
	bh.Start()
}
