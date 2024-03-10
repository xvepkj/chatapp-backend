package main

import (
	"github.com/gin-gonic/gin"
	"github.com/xvepkj/chatapp-backend/handlers"
	"github.com/xvepkj/chatapp-backend/models"
	"github.com/xvepkj/chatapp-backend/utils"
)

func main() {
	db, err := utils.ConnectDB()
	if err != nil {
		panic("failed to connect to database")
	}

	sqlDB, _ := db.DB()

	defer sqlDB.Close()

	err = db.AutoMigrate(&models.User{}, &models.Message{})
	if err != nil {
		panic("failed to migrate database")
	}

	router := gin.Default()

	router.POST("/users", func(c *gin.Context) {
		handlers.CreateUser(c, db)
	})

	router.GET("/users/:id", func(c *gin.Context) {
		handlers.GetUserByID(c, db)
	})

	router.POST("/messages", func(c *gin.Context) {
		handlers.AddMessage(c, db)
	})

	router.GET("/messages/:senderID/:receiverID", func(c *gin.Context) {
		handlers.GetMessagesBetween(c, db)
	})

	router.Run()
}
