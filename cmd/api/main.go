package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"database/sql"
	"peapong-auth/internal/auth/api"
	"peapong-auth/internal/auth/repository"
	"peapong-auth/internal/auth/service"
	mydb "peapong-auth/pkg/database"

	"github.com/joho/godotenv"
)

type App struct {
	DB          *sql.DB
	Server      *http.Server
	AuthHandler *api.AuthHandler
}

func NewApp(port string) *App {
	//Connect DB
	db := mydb.ConnectDB()

	// Inject service
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	authHandler := api.NewAuthHandler(authService)

	//Router
	mux := http.NewServeMux()

	mux.HandleFunc("POST /signup", authHandler.RegisterNewUserHandler)
	mux.HandleFunc("GET /user", authHandler.FindUserNameByIdHandler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	return &App{
		DB:          db,
		Server:      srv,
		AuthHandler: authHandler,
	}
}

func (a *App) Run() {

	go func() {
		fmt.Printf("Starting server at %s...\n", a.Server.Addr)
		if err := a.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	a.Shutdown()
}

func (a *App) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.Server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Closing database connection...")
	if err := a.DB.Close(); err != nil {
		log.Printf("Error closing db: %v", err)
	}

	log.Println("Server exited properly")
}

func main() {
	// 1. Load Config
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// 2. Initialize App
	app := NewApp(port)

	// 3. Start App
	app.Run()
}
