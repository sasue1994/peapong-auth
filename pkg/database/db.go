package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type dbConfig struct {
	db *sql.DB
}

// ConnectDB สร้างและคืนค่าการเชื่อมต่อฐานข้อมูล
// TODO: ควรย้าย Connection String ไปใช้ Environment Variables แทนการ Hardcode
func ConnectDB() *sql.DB {
	connStr := "host=localhost port=5432 user=postgres password=kAZ3aKOsAdPs37Rr  dbname=postgres sslmode=disable"

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db
}
