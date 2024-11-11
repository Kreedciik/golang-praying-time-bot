package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Timings struct {
	Fajr    string
	Sunrise string
	Dhuhr   string
	Asr     string
	Sunset  string
	Maghrib string
	Isha    string
}
type ResponseData struct {
	Timings Timings `json:"timings"`
	Date    Date    `json:"date"`
	Meta    Meta    `json:"meta"`
}
type Hijri struct {
	Date string `json:"date"`
}

type T struct {
	Day     string `json:"day"`
	Weekday struct {
		En string `json:"en"`
	} `json:"weekday"`
	Month struct {
		Number int    `json:"number"`
		En     string `json:"en"`
	} `json:"month"`
	Year string `json:"year"`
}

type Gregorian struct {
	Date    string `json:"date"`
	Day     string `json:"day"`
	Weekday struct {
		En string `json:"en"`
	} `json:"weekday"`
	Month struct {
		Number int    `json:"number"`
		En     string `json:"en"`
	} `json:"month"`
	Year string `json:"year"`
}
type Date struct {
	Hijri     Hijri     `json:"hijri"`
	Gregorian Gregorian `json:"gregorian"`
}
type Method struct {
	Name string `json:"name"`
}
type Meta struct {
	Timezone string `json:"timezone"`
	Method   Method `json:"method"`
}
type Response struct {
	Data ResponseData `json:"data"`
}

type PrayingSchedule struct {
	Timing        Timings
	GregorianDate string
	HijriDate     string
	CurrentMonth  string
	CurrentDay    string
}

func getPrayingTime(latitude, longitude float64, date string) Response {
	var response Response = Response{}
	url := fmt.Sprintf("https://api.aladhan.com/v1/timings/%v?method=3&latitude=%f&longitude=%f", date, latitude, longitude)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Error while getting data", err)
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Fatal("Error parsing response:", err)
	}
	return response
}
