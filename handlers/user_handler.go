package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-whatsapp-api/models"
	"go-whatsapp-api/services"
	"log"
	"net/http"
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
	log.Printf("Creating user")

	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Received user data: %+v", newUser)
	if !whatsAppService.ValidatePhone(newUser.WhatsApp.Phone) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid WhatsApp number"})
		return
	}

	if err := models.CreateUser(newUser); err != nil {
		if err.Error() == fmt.Sprintf("phone number %s is already registered", newUser.WhatsApp.Phone) {
			c.JSON(http.StatusConflict, gin.H{"error": "Phone number already registered"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if newUser.WhatsApp.OptIn {
		msg := fmt.Sprintf("Hello %s! Thanks for registering.", newUser.Username)
		if _, err := whatsAppService.SendMessage(newUser.WhatsApp.Phone+"@s.whatsapp.net", msg); err != nil {
			c.JSON(http.StatusCreated, gin.H{
				"user":    newUser,
				"warning": "WhatsApp message failed: " + err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusCreated, newUser)
}

func UpdateUser(c *gin.Context) {
	userid := c.Query("userid")
	if userid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	var updatedUser models.User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingUser, err := models.GetUserByID(userid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	updatedUser.WhatsApp.Phone = existingUser.WhatsApp.Phone

	err = models.UpdateUser(userid, updatedUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func GetUser(c *gin.Context) {
	userid := c.Query("userid")
	if userid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	user, err := models.GetUserByID(userid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func LoginUser(c *gin.Context) {
	userid := c.Query("userid")
	if userid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	user, err := models.GetUserByID(userid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"userid": user.ID,
		"role":   user.Role,
	})
}

func UpdateNews(c *gin.Context) {
	var news models.News
	if err := c.ShouldBindJSON(&news); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := models.UpdateNews(news)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "News updated successfully"})
}

func GetNews(c *gin.Context) {
	news, err := models.GetNews()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, news)
}

func UpdateCompetition(c *gin.Context) {
	var competition models.Competition
	if err := c.ShouldBindJSON(&competition); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := models.UpdateCompetition(competition)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Competition info updated successfully"})
}

func GetCompetition(c *gin.Context) {
	competition, err := models.GetCompetition()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, competition)
}

func UpdateSchedule(c *gin.Context) {
	var schedule models.Schedule
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := models.UpdateSchedule(schedule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Schedule updated successfully"})
}

func GetSchedule(c *gin.Context) {
	schedule, err := models.GetSchedule()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schedule)
}
