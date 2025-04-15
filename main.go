package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"wechatRobot/app/config"
	"wechatRobot/app/sdk"
	"wechatRobot/app/service"
	"wechatRobot/app/storage"
)

func main() {
	// Initialize WeChat SDK
	err := sdk.Initialize(config.AppConfig.Debug)
	if err != nil {
		log.Fatalf("Failed to initialize SDK: %v", err)
	}
	defer sdk.Cleanup()

	// Initialize storage and services
	db := storage.NewDatabase()
	wxService := service.NewWeChatService(db)

	// Wait for login
	if err = wxService.WaitForLogin(); err != nil {
		log.Fatalf("Login failed: %v", err)
	}

	// Initialize contacts
	go wxService.InitializeContacts()

	// Start message listener
	log.Printf("Starting message listener for user: %s", wxService.GetCurrentUser())
	wxService.StartMessageListener()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("Shutting down gracefully...")
}
