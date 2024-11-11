package main

import (
	"database/sql"
	"fmt"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func initialScreen(update telego.Update, bot *telego.Bot, db *sql.DB, mainKeyboard *telego.ReplyKeyboardMarkup) {
	userID := update.Message.Chat.ID
	firstName := update.Message.From.FirstName
	latitude := update.Message.Location.Latitude
	longitude := update.Message.Location.Longitude
	response := getPrayingTime(latitude, longitude, "{}")
	newUser := User{user_id: userID, first_name: firstName, time_zone: response.Data.Meta.Timezone, location: Location{latitude, longitude}}
	// Inserting new user
	insertNewUser(db, newUser)
	// Inserting new schedule based on user's location
	prayingSchedule := response
	insertPrayingSchedule(db, userID, prayingSchedule)
	_, _ = bot.SendMessage(tu.Message(tu.ID(update.Message.Chat.ID), fmt.Sprintf("Your timezone: %v", prayingSchedule.Data.Meta.Timezone)).WithReplyMarkup(mainKeyboard))
}

func welcomeScreen(bot *telego.Bot, update telego.Update, startKeyboard *telego.ReplyKeyboardMarkup) {
	_, _ = bot.SendMessage(tu.Message(tu.ID(update.Message.Chat.ID), fmt.Sprintf("Assalamu-alekum %s!\nWelcome to this bot!\nThrough this bot you can get prayer timings based on your location.", update.Message.From.FirstName)).WithReplyMarkup(startKeyboard))
	_, _ = bot.SendMessage(tu.Message(tu.ID(update.Message.Chat.ID), "Please, specify your location").WithReplyMarkup(startKeyboard))
}
