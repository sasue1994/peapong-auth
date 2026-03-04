package repository

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// UserRepository จัดการกับการเชื่อมต่อฐานข้อมูลและการดำเนินการเกี่ยวกับผู้ใช้
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository สร้างอินสแตนซ์ใหม่ของ UserRepository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// RegisterNewUser เพิ่มผู้ใช้ใหม่ลงในฐานข้อมูล
func (r *UserRepository) RegisterNewUser(shopID int, username string, password string) bool {
	query := `INSERT INTO users (shop_id, username, password_hash) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(query, shopID, username, getHashPassword(password))
	if err != nil {
		log.Printf("Insert Error: %v", err)
		return false
	}
	return true
}

// FindUserNameByID ค้นหาชื่อผู้ใช้ด้วย ID
func (r *UserRepository) FindUserNameByID(userID string) (string, error) {
	query := `SELECT username FROM users WHERE id = $1`
	var username string
	err := r.db.QueryRow(query, userID).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil // ไม่พบผู้ใช้, ไม่ใช่ error
		}
		log.Printf("Query Error: %v", err)
		return "", err
	}
	return username, nil
}

// getHashPassword สร้าง HashPassword
func getHashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}
