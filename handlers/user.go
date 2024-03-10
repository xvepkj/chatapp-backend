package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xvepkj/chatapp-backend/db"
	"github.com/xvepkj/chatapp-backend/models"
	"gorm.io/gorm"
)

func CreateUser(c *gin.Context, dbConn *gorm.DB) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.CreateUser(dbConn, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user created successfully", "user": user})
}

func GetUserByID(c *gin.Context, dbConn *gorm.DB) {
	userID := c.Param("id")

	user, err := db.GetUserByUsername(dbConn, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func UpdateUser(c *gin.Context, dbConn *gorm.DB) {
	var user models.User
	username := c.Param("username") // Assuming the parameter name is "username"

	if _, err := db.GetUserByUsername(dbConn, username); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Bind the updated user data from the request body
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the user in the database
	if err := db.UpdateUser(dbConn, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with success message
	c.JSON(http.StatusOK, gin.H{"message": "user updated successfully", "user": user})
}

func DeleteUser(c *gin.Context, dbConn *gorm.DB) {
	userID := c.Param("id")

	if err := db.DeleteUser(dbConn, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}
