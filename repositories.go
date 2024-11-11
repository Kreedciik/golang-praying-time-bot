package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
)

type Location struct{ latitude, longitude float64 }
type User struct {
	user_id    int64
	first_name string
	time_zone  string
	location   Location
}

func createUserTable(db *sql.DB) {
	_, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS users (
    			user_id bigint PRIMARY KEY,
				first_name varchar NOT NULL,
				time_zone varchar NOT NULL,
				latitude float NOT NULL,
				longitude float NOT NULL
)`)
	if err != nil {
		log.Fatal("Error while creating users table", err)
	}
}

func insertNewUser(db *sql.DB, user User) {
	_, err := db.Exec(
		`INSERT INTO users (user_id, first_name, time_zone, latitude, longitude) VALUES ($1, $2, $3, $4, $5)`, user.user_id, user.first_name, user.time_zone, user.location.longitude, user.location.latitude)
	if err != nil {
		fmt.Println("Error while inserting user into table", err)
		panic(err)
	}
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()
}

func createPrayingTimeTable(db *sql.DB) {
	_, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS praying_time (
    			pray_time_id bigint,
				timing jsonb,
				gregorian_date date,
				hijri_date date,
				current_month varchar,
				current_week varchar,
				calc_method varchar,
				user_id bigint,
				CONSTRAINT fk_user_id FOREIGN KEY(user_id) REFERENCES users(user_id)
)`)
	if err != nil {
		log.Fatal("Error while creating praying time table", err)
	}
}

func insertPrayingSchedule(db *sql.DB, userID int64, prayingSchedule Response) {
	timings := prayingSchedule.Data.Timings
	jsonSchedule, err := json.Marshal(map[string]string{
		"Fajr":    timings.Fajr,
		"Sunrise": timings.Sunrise,
		"Dhuhr":   timings.Dhuhr,
		"Asr":     timings.Asr,
		"Maghrib": timings.Maghrib,
		"Isha":    timings.Isha,
	})
	if err != nil {
		fmt.Println(err)
	}
	prayTimeID := userID
	gregorianDate := prayingSchedule.Data.Date.Gregorian.Date
	hijriDate := prayingSchedule.Data.Date.Hijri.Date
	currentMonth := prayingSchedule.Data.Date.Gregorian.Month.En
	currentWeek := prayingSchedule.Data.Date.Gregorian.Weekday.En
	calcMethod := prayingSchedule.Data.Meta.Method.Name
	_, err2 := db.Exec(
		`INSERT INTO praying_time (pray_time_id, timing, gregorian_date, hijri_date, current_month, current_week, calc_method, user_id) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`, prayTimeID, jsonSchedule, gregorianDate, hijriDate, currentMonth, currentWeek, calcMethod, userID)
	if err2 != nil {
		log.Fatal("Error while inserting praying time", err)
	}

}

func getUser(db *sql.DB, userID int64) (Location, string, error) {
	var userFirstName string
	var location Location
	err := db.QueryRow(`SELECT first_name, latitude, longitude FROM users WHERE user_id = $1`, userID).Scan(&userFirstName, &location.latitude, &location.longitude)
	if err != nil {
		fmt.Println("Error while getting user", err)
		return location, "", err
	}
	recover()
	return location, userFirstName, nil
}

func getPrayingSchedule(db *sql.DB, userID int64) (PrayingSchedule, error) {
	var prayingSchedule PrayingSchedule
	var schedule string
	err := db.QueryRow(`SELECT timing, gregorian_date, hijri_date, current_month, current_week FROM praying_time WHERE pray_time_id = $1`, userID).Scan(&schedule, &prayingSchedule.GregorianDate, &prayingSchedule.HijriDate, &prayingSchedule.CurrentMonth, &prayingSchedule.CurrentDay)
	if err != nil {
		log.Fatal("Error while getting praying schedule", err)
		return prayingSchedule, err
	}
	err2 := json.Unmarshal([]byte(schedule), &prayingSchedule.Timing)
	if err2 != nil {
		log.Fatal("Error while unmarshalling praying schedule", err2)
	}
	return prayingSchedule, nil
}
