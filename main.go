package main

import (
	"encoding/json"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/xvepkj/chatapp-backend/handlers"
	"github.com/xvepkj/chatapp-backend/models"
	"github.com/xvepkj/chatapp-backend/utils"
	"gorm.io/gorm"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Map to store WebSocket connections of users
var userConnections = make(map[string]*websocket.Conn)

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

	router.POST("/users/register", func(c *gin.Context) {
		handlers.CreateUser(c, db)
	})

	router.POST("/users/login", func(c *gin.Context) {
		handlers.GetUser(c, db)
	})

	router.GET("/users/:id", authMiddleware, func(c *gin.Context) {
		handlers.GetUserByID(c, db)
	})

	router.POST("/messages", authMiddleware, func(c *gin.Context) {
		handlers.AddMessage(c, db)
	})

	router.GET("/messages/:senderID/:receiverID", authMiddleware, func(c *gin.Context) {
		handlers.GetMessagesBetween(c, db)
	})

	router.GET("/ws", func(c *gin.Context) {
		handleWebSocketConnection(c, db)
	})

	router.Run()
}

// authMiddleware is a middleware function to authenticate JWT tokens
func authMiddleware(c *gin.Context) {
	// Get the JWT token from the Authorization header
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization token"})
		c.Abort()
		return
	}

	// Parse the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("signing-key"), nil // Replace "signing-key" with your actual secret key
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		c.Abort()
		return
	}

	// Extract user information from the token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
		c.Abort()
		return
	}

	// Add user information to the Gin context
	username, ok := claims["authenticated_user"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username in token"})
		c.Abort()
		return
	}
	c.Set("authenticated_user", username)

	// Continue to the next middleware or route handler
	c.Next()
}

func sendWebSocketMessage(message models.Message) {
	// Find recipient's WebSocket connection
	recipientConn := userConnections[message.ReceipientID]
	if recipientConn != nil {
		// Convert message to JSON
		messageJSON, err := json.Marshal(message)
		if err != nil {
			// Handle error
			return
		}

		// Send message to recipient
		err = recipientConn.WriteMessage(websocket.TextMessage, messageJSON)
		if err != nil {
			// Handle error
			return
		}
	} else {
		// Handle case where recipient is not found
		return
	}
}

func handleWebSocketConnection(c *gin.Context, db *gorm.DB) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		// Handle error
		return
	}
	defer conn.Close()

	// Read messages from WebSocket connection
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			// Handle error
			break
		}

		// Process received message
		var receivedMessage models.Message
		err = json.Unmarshal(p, &receivedMessage)
		if err != nil {
			// Handle error
			break
		}

		handlers.AddMessageWebSocket(db, &receivedMessage)

		// Add user connection to map if not already present
		if _, ok := userConnections[receivedMessage.ReceipientID]; !ok {
			userConnections[receivedMessage.ReceipientID] = conn
		}

		// Broadcast message to recipient
		sendWebSocketMessage(receivedMessage)
	}
}
