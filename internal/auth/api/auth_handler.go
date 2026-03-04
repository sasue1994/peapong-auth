package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"peapong-auth/internal/auth/service"
	"peapong-auth/models"
	"peapong-auth/pkg/security"
	"strings"
)

// Inject service
type AuthHandler struct {
	authService *service.AuthService
}

// ประกาศ Constructor
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// ประกาศ Method
// RegisterNewUserHandler
func (h *AuthHandler) RegisterNewUserHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	apiKey := req.Header.Get("x-api-key")

	isKeyValid := security.ValidationAPIKey(apiKey)

	if isKeyValid == false {
		http.Error(res, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var user models.User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		http.Error(res, "Bad Request: Invalid JSON", http.StatusBadRequest)
		return
	}

	if user.ShopID == 0 || strings.TrimSpace(user.Username) == "" || strings.TrimSpace(user.Password) == "" {
		http.Error(res, "Bad Request: Missing required fields", http.StatusBadRequest)
		return
	}

	if h.authService.RegisterNewUser(user.ShopID, user.Username, user.Password) {
		res.WriteHeader(http.StatusCreated)
	} else {
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *AuthHandler) FindUserNameByIdHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	apiKey := req.Header.Get("x-api-key")

	fmt.Println(apiKey)

	isKeyValid := security.ValidationAPIKey(apiKey)

	fmt.Println(isKeyValid)

	if isKeyValid == false {
		http.Error(res, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userId := req.URL.Query().Get("userId")
	if userId == "" {
		http.Error(res, "Bad Request: Missing required fields", http.StatusBadRequest)
		return
	}
	username, err := h.authService.FindUserNameByID(userId)

	if err != nil {
		fmt.Println(err)
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	userName := username
	if userName == "" {
		http.Error(res, "Not Found", http.StatusNotFound)
		return
	}
	fmt.Fprintf(res, "User Name: %s", userName)

	data := map[string]string{"username": userName}

	res.Header().Set("Content-Type", "application/json")

	res.WriteHeader(http.StatusOK)

	err = json.NewEncoder(res).Encode(data)

	if err != nil {
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}
