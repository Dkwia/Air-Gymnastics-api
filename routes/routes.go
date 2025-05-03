package routes

import (
	"github.com/gin-gonic/gin"
	"go-whatsapp-api/handlers"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		users := api.Group("/users")
		{
			users.GET("/", handlers.GetAllWhatsAppUsers)
			users.POST("/", handlers.CreateUser)
		}
	}
}
