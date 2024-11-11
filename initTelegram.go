package main

import (
	"database/sql"
	"fmt"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"os"
	"time"
)

func sendScheduleToUser(bot *telego.Bot, userID int64, schedule PrayingSchedule) {
	parsedGregorianDate, _ := time.Parse(time.RFC3339, schedule.GregorianDate)
	parsedHijriDate, _ := time.Parse(time.RFC3339, schedule.HijriDate)
	messageText := fmt.Sprintf(
		"☪️Praying schedule:\n"+
			"\n"+
			"📆 %v\n"+
			"📆 %v\n"+
			"➖➖➖➖➖➖➖➖\n"+
			"⏰ %v | Fajr\n"+
			"⏰ %v | Sunrise\n"+
			"⏰ %v | Dhuhr\n"+
			"⏰ %v | Asr\n"+
			"⏰ %v | Maghrib\n"+
			"⏰ %v | Isha\n"+
			"➖➖➖➖➖➖➖➖",
		fmt.Sprintf("%v %v, %v", schedule.CurrentMonth, schedule.CurrentDay, parsedGregorianDate.Year()), // Example: "14-May 2023"
		fmt.Sprintf("%v (Hijri date)", parsedHijriDate.Format("02.01.2006")),                             // Example: "23-shavvol 1444"
		schedule.Timing.Fajr,
		schedule.Timing.Sunrise,
		schedule.Timing.Dhuhr,
		schedule.Timing.Asr,
		schedule.Timing.Maghrib,
		schedule.Timing.Isha,
	)
	_, _ = bot.SendMessage(tu.Message(tu.ID(userID), messageText))
}

func initTelegramBot(db *sql.DB) {
	botToken := os.Getenv("BOT_TOKEN")
	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	updates, _ := bot.UpdatesViaLongPolling(nil)
	defer bot.StopLongPolling()

	// Create a main keyboard
	mainKeyboard := tu.Keyboard(
		tu.KeyboardRow(
			tu.KeyboardButton("Praying schedule for today"),
		),
		tu.KeyboardRow(
			tu.KeyboardButton("Praying schedule for tomorrow")),
	).WithResizeKeyboard().WithInputFieldPlaceholder("Select something...")
	// Create a bot handler
	bh, _ := th.NewBotHandler(bot, updates)
	defer bh.Stop()

	// Handle /start command
	bh.Handle(startCommandHandler, th.CommandEqual("start"))

	// Handle "Praying schedule for today" button
	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		userID := update.Message.Chat.ID
		schedule, err := getPrayingSchedule(db, userID)
		if err != nil {
			fmt.Println("ERROR ", err)
		}
		sendScheduleToUser(bot, userID, schedule)
	}, th.TextEqual("Praying schedule for today"))

	// Handle "Praying schedule for tomorrow" button
	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		userID := update.Message.Chat.ID
		tomorrow := time.Now().AddDate(0, 0, 1)
		location, _, err := getUser(db, userID)
		if err != nil {
			fmt.Println("ERROR ", err)
		}
		response := getPrayingTime(location.latitude, location.longitude, tomorrow.Format("02-01-2006"))
		var prayingSchedule PrayingSchedule
		gregorianDate, _ := time.Parse("02-01-2006", response.Data.Date.Gregorian.Date)
		hijriDate, _ := time.Parse("02-01-2006", response.Data.Date.Hijri.Date)
		prayingSchedule.CurrentDay = response.Data.Date.Gregorian.Day
		prayingSchedule.CurrentMonth = response.Data.Date.Gregorian.Month.En
		prayingSchedule.GregorianDate = gregorianDate.Format(time.RFC3339)
		prayingSchedule.HijriDate = hijriDate.Format(time.RFC3339)
		prayingSchedule.Timing = response.Data.Timings
		sendScheduleToUser(bot, userID, prayingSchedule)
	}, th.TextEqual("Praying schedule for tomorrow"))

	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		// Define user in a database
		userID := update.Message.Chat.ID
		fmt.Println("ANY MESSAGE")
		if update.Message.Location != nil {
			// Handle "Send my location" button
			initialScreen(update, bot, db, mainKeyboard)
		} else {
			_, _ = bot.SendMessage(tu.Message(tu.ID(userID), "Thank you!").WithReplyMarkup(mainKeyboard))
		}
	}, th.AnyMessage())

	// Start handling updates
	bh.Start()
}
