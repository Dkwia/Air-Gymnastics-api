package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"
	
	_ "github.com/mattn/go-sqlite3"
	"go-whatsapp-api/handlers"
	"go-whatsapp-api/models"
	"go-whatsapp-api/routes"
	"go-whatsapp-api/services"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := models.InitDB(); err != nil {
		log.Fatalf("Failed to init DB: %v", err)
	}
	defer models.DB.Close()

	ws, err := services.NewWhatsAppService()
	if err != nil {
		log.Fatalf("Failed to init WhatsApp: %v", err)
	}
	defer ws.Disconnect()

	handlers.InitHandlers(ws)

	r := gin.Default()
	routes.SetupRoutes(r)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := r.Run(":8080"); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down...")
	time.Sleep(1 * time.Second)
}
