package handlers

import (
	"net/http"
	"go-gin-api/models"
	"go-gin-api/services"
	"github.com/gin-gonic/gin"
)

var whatsAppService *services.WhatsAppService

func init() {
	var err error
	whatsAppService, err = services.NewWhatsAppService()
	if err != nil {
		panic("Failed to initialize WhatsApp service: " + err.Error())
	}
}

func GetUsers(c *gin.Context) {
	c.JSON(http.StatusOK, models.Users)
}

func GetUser(c *gin.Context) {
	id := c.Param("id")

	for _, user := range models.Users {
		if user.ID == id {
			c.JSON(http.StatusOK, user)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
}

func CreateUser(c *gin.Context) {
	var newUser models.User

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !whatsAppService.ValidatePhone(newUser.WhatsApp.Phone) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid WhatsApp phone number"})
		return
	}

	models.Users = append(models.Users, newUser)
	
	if newUser.WhatsApp.OptIn {
		msg := fmt.Sprintf("Привет, %s! Спасибо за регистрацию.", newUser.Username)
		msgID, err := whatsAppService.SendMessage(newUser.WhatsApp.Phone, msg)
		
		if err != nil {
			c.JSON(http.StatusCreated, gin.H{
				"user": newUser,
				"warning": "User created but WhatsApp message failed: " + err.Error(),
			})
			return
		}
		
		for i, u := range models.Users {
			if u.ID == newUser.ID {
				models.Users[i].WhatsApp.LastMsgID = msgID
				break
			}
		}
	}

	c.JSON(http.StatusCreated, newUser)
}

func SendWhatsAppMessage(c *gin.Context) {
	id := c.Param("id")
	var request struct {
		Message string `json:"message"`
	}
	
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user *models.User
	for _, u := range models.Users {
		if u.ID == id {
			user = &u
			break
		}
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	if !user.WhatsApp.OptIn {
		c.JSON(http.StatusForbidden, gin.H{"error": "User has not opted in for WhatsApp messages"})
		return
	}

	msgID, err := whatsAppService.SendMessage(user.WhatsApp.Phone, request.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i, u := range models.Users {
		if u.ID == id {
			models.Users[i].WhatsApp.LastMsgID = msgID
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "WhatsApp message sent successfully",
		"message_id": msgID,
	})
}
