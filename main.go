package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file ", err)
		os.Exit(1)
	}
	db := createDBInstance()
	createUserTable(db)
	createPrayingTimeTable(db)
	defer initTelegramBot(db)
}
