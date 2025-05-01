package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"go-gin-api/routes"
	"go-gin-api/services"
	"github.com/gin-gonic/gin"
)

func main() {
	whatsAppService, err := services.NewWhatsAppService()
	if err != nil {
		log.Fatalf("Failed to initialize WhatsApp service: %v", err)
	}
	defer whatsAppService.Disconnect()

	r := gin.Default()

	routes.SetupRoutes(r)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := r.Run(":8080"); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down gracefully...")
}
