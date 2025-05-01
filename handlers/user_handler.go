package handlers

import (
	"net/http"
	"go-gin-api/models"
	"github.com/gin-gonic/gin"
)

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

	models.Users = append(models.Users, newUser)
	c.JSON(http.StatusCreated, newUser)
}
