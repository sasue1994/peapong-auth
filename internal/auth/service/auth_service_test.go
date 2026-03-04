package service_test

import (
	"testing"

	"peapong-auth/internal/auth/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// สร้าง Mock Repo
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindUserNameByID(userID string) (string, error) {
	args := m.Mock.Called(userID)
	// คืนค่าตามที่ setup ไว้ (ค่าที่ 0 คือ string, ค่าที่ 1 คือ error)
	return args.String(0), args.Error(1)
}

func (m *MockUserRepository) RegisterNewUser(shopID int, username string, password string) bool {
	args := m.Mock.Called(shopID, username, password)
	return args.Bool(0)
}

func TestHashPassword(t *testing.T) {

	mockRepo := new(MockUserRepository)
	authService := service.NewAuthService(mockRepo)
	password := "secret123"

	hash, err := authService.HashPassword(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)
}

func TestCheckPasswordHash(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := service.NewAuthService(mockRepo)
	password := "secret123"
	hash, _ := authService.HashPassword(password)

	// Case 1: Password ตรงกัน
	match := authService.CheckPasswordHash(password, hash)
	assert.True(t, match, "Password should match")

	// Case 2: Password ไม่ตรงกัน
	mismatch := authService.CheckPasswordHash("wrongpassword", hash)
	assert.False(t, mismatch, "Password should not match")
}
