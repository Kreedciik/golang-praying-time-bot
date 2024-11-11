package main

import (
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func startCommandHandler(bot *telego.Bot, update telego.Update) {
	// Create a start keyboard
	startKeyboard := tu.Keyboard(
		tu.KeyboardRow(
			tu.KeyboardButton("Send my location").WithRequestLocation(),
		),
	).WithResizeKeyboard().WithInputFieldPlaceholder("Select something...")
	
	welcomeScreen(bot, update, startKeyboard)
}
