package handlers

import (
	"fmt"
	"net/http"
	"go-whatsapp-api/models"
	"go-whatsapp-api/services"
	"github.com/gin-gonic/gin"
)

var whatsAppService *services.WhatsAppService

func InitHandlers(ws *services.WhatsAppService) {
	whatsAppService = ws
}

func GetAllWhatsAppUsers(c *gin.Context) {
	users, err := models.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func CreateUser(c *gin.Context) {
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !whatsAppService.ValidatePhone(newUser.WhatsApp.Phone) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid WhatsApp number"})
		return
	}

	if err := models.CreateUser(newUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if newUser.WhatsApp.OptIn {
		msg := fmt.Sprintf("Hello %s! Thanks for registering.", newUser.Username)
		if _, err := whatsAppService.SendMessage(newUser.WhatsApp.Phone, msg); err != nil {
			c.JSON(http.StatusCreated, gin.H{
				"user": newUser,
				"warning": "WhatsApp message failed: " + err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusCreated, newUser)
}
