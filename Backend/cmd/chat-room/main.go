package main

import (
	// "fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/akshayjha21/Chat-App-in-GO/Backend/internal/chat"
	handler "github.com/akshayjha21/Chat-App-in-GO/Backend/internal/chat/Handler"
	"github.com/akshayjha21/Chat-App-in-GO/Backend/internal/config"
	"github.com/akshayjha21/Chat-App-in-GO/Backend/internal/storage/postgres"
	"github.com/gofiber/fiber/v2"
)

func main() {

	cfg := config.MustLoad()
	//db connection

	log.Printf("Conecting to the database")
	db, err := postgres.New(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// userHandler := &user.Handler{DB: db}
	userHandler := &handler.UserHandler{DB: db}
	chatHandler := &handler.Chathandler{DB: db}
	//1. hub ko initialize krenge
	hub := chat.NewHub(db)
	// Get the generic sql.DB object from GORM to close it
	//2. hub ko background me execute krenge
	go hub.Run()
	
	app:=fiber.New();
	// 3. Route Handler Update karenge
	// Note: Hum seedha chat.ServeWs nahi de sakte kyunki usse 'hub' chahiye.
	// Isliye hum ek "Closure" (anonymous function) use karenge.

	app.Post("/register",userHandler.Registerhandler)
	app.Post("/login",userHandler.LoginHandler)
	app.Post("/createRoom",chatHandler.CreateChatRoom)
	app.Post("/joinRoom",chatHandler.JoinRoom)

	go func() {
		log.Printf("Server starting on %s", cfg.Addr)
		if err := app.Listen(cfg.Addr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
		// Graceful Shutdown logic
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
	log.Println("Server stopping.....")
	if sqlDB, err := db.Db.DB(); err == nil {
		log.Println("Closing database connection...")
		sqlDB.Close()
	}
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited properly")
}
