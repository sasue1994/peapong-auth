package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"golang.org/x/crypto/bcrypt"
)

// Global
var db *sql.DB

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func GetHashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type Account struct {
	Username string `json:"username"`
	Password string `json:"password"`
	ShopID   int    `json:"shop_id"`
}

func signup(res http.ResponseWriter, req *http.Request) {
	var account Account

	fmt.Printf("Method: %s", req.Method)

	apiKey := req.Header.Get("x-api-key")

	if apiKey != "xYZ1234" {
		http.Error(res, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := json.NewDecoder(req.Body).Decode(&account)
	if err != nil {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	shopID := account.ShopID
	username := account.Username
	password := account.Password

	if shopID == 0 || username == "" || password == "" {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	if len(strings.TrimSpace(username)) == 0 || len(strings.TrimSpace(password)) == 0 {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	if RegisterNewUser(shopID, username, password) {
		res.WriteHeader(http.StatusCreated)
	}

}

func RegisterNewUser(shopID int, username string, password string) bool {

	query := `INSERT INTO users (shop_id, username, password_hash) VALUES ($1, $2, $3)`
	_, err := db.Exec(query, shopID, username, GetHashPassword(password)) // ใช้ GetHashPassword แทน encrypetionPassword)
	if err != nil {
		log.Printf("Insert Error: %v", err)
		return false
	}
	return true
}

func login(respons http.ResponseWriter, request *http.Request) {}

func DatabaseConnection() *sql.DB {

	connStr := "host=localhost port=5432 user=postgres password=kAZ3aKOsAdPs37Rr  dbname=postgres sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}

	db.SetMaxOpenConns(5)                  // จำนวน connection สูงสุดที่เปิดได้
	db.SetMaxIdleConns(5)                  // จำนวน connection ที่เปิดค้างไว้รอใช้งาน
	db.SetConnMaxLifetime(5 * time.Minute) // อายุสูงสุดของแต่ละ connection

	err = db.Ping()
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	fmt.Println("Successfully connected to database!")
	return db
}

func findUserNameById(userId string) string {

	query := `SELECT username FROM users WHERE id = $1`
	var username string
	err := db.QueryRow(query, userId).Scan(&username)
	if err != nil {
		log.Printf("Query Error: %v", err)
	}
	return username

}

func getUserName(res http.ResponseWriter, req *http.Request) {

	fmt.Printf("method : %s", req.Method)

	headerString := req.Header.Get("Content-Type")

	fmt.Printf("header : %s", headerString)

	userId := req.URL.Query().Get("userId")

	if userId == "" {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	userName := findUserNameById(userId)
	fmt.Fprintf(res, "User Name: %s", userName)

	if userName == "" {
		http.Error(res, "Not Found", http.StatusNotFound)
		return
	}

	data := map[string]string{"username": userName}

	res.Header().Set("Content-Type", "application/json")

	res.WriteHeader(http.StatusOK)

	err := json.NewEncoder(res).Encode(data)

	if err != nil {
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}

func main() {
	db = DatabaseConnection()
	// defer db.Close()

	defer func() {

		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	// http.HandleFunc("/getUserId", getUserName)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /getUserId", getUserName)
	mux.HandleFunc("POST /signup", signup)

	port := ":8081"
	err := http.ListenAndServe(port, mux)

	fmt.Printf("Starting server at port %s...\n", port)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
