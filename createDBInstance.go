package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func createDBInstance() *sql.DB {
	pgPassword := os.Getenv("PG_PASSWORD")
	pgDbName := os.Getenv("PG_DB")
	pgUser := os.Getenv("PG_USER")
	connStr := fmt.Sprintf("user=%v password=%v dbname=%v sslmode=disable", pgUser, pgPassword, pgDbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	log.Println("Successfully connected to DB")
	return db
}
