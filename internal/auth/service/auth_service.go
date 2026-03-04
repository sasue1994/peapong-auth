package service

import (
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	FindUserNameByID(userID string) (string, error)
	RegisterNewUser(shopID int, username string, password string) bool
}

type AuthService struct {
	userRepo UserRepository
}

// ประกาศ Constructor
func NewAuthService(userRepo UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *AuthService) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// FindUserNameByID ค้นหาชื่อผู้ใช้ด้วย ID
func (s *AuthService) FindUserNameByID(userID string) (string, error) {
	return s.userRepo.FindUserNameByID(userID)
}

// RegisterNewUser เพิ่มผู้ใช้ใหม่ลงในฐานข้อมูล
func (s *AuthService) RegisterNewUser(shopID int, username string, password string) bool {
	return s.userRepo.RegisterNewUser(shopID, username, password)
}
