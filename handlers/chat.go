package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xvepkj/chatapp-backend/db"
	"github.com/xvepkj/chatapp-backend/models"
	"gorm.io/gorm"
)

func AddMessage(c *gin.Context, dbConn *gorm.DB) {
	var message models.Message
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.AddMessage(dbConn, &message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "message added successfully", "user": message})
}

func GetMessagesBetween(c *gin.Context, dbConn *gorm.DB) {
	senderID := c.Param("senderID")
	receiverId := c.Param("receiverID")

	messages, err := db.GetMessagesBetween(dbConn, senderID, receiverId)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "messages not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}
