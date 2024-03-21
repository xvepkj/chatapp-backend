package main

import (
	"encoding/json"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/juju/ratelimit"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	docs "github.com/xvepkj/chatapp-backend/docs"
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

// Initialize a rate limiter with a maximum of 50 requests per minute
var limiter = ratelimit.NewBucketWithRate(60, 50)

// RateLimitMiddleware applies rate limiting to incoming requests
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the client's request rate exceeds the limit
		if limiter.TakeAvailable(1) == 0 {
			// Reject the request with a 429 Too Many Requests status code
			c.AbortWithStatusJSON(429, gin.H{"error": "too many requests"})
			return
		}
		// Continue to the next middleware or route handler
		c.Next()
	}
}

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

	docs.SwaggerInfo.BasePath = "/"

	router := gin.Default()

	router.Use(RateLimitMiddleware())

	router.POST("/users/register", func(c *gin.Context) {
		handlers.CreateUser(c, db)
	})

	router.POST("/users/login", func(c *gin.Context) {
		handlers.GetUser(c, db)
	})

	router.GET("/users/:id", authMiddleware, func(c *gin.Context) {
		handlers.GetUserByID(c, db)
	})

	router.GET("/users", authMiddleware, func(c *gin.Context) {
		handlers.GetAllUsernames(c, db)
	})

	router.POST("/messages", authMiddleware, func(c *gin.Context) {
		handlers.AddMessage(c, db)
	})

	router.GET("/messages/:senderID/:receiverID", authMiddleware, func(c *gin.Context) {
		handlers.GetMessagesBetween(c, db)
	})

	router.GET("export/messages/:senderID/:receiverID", authMiddleware, func(c *gin.Context) {
		handlers.ExportMessagesToExcel(c, db)
	})

	router.GET("/ws", authMiddleware, func(c *gin.Context) {
		handleWebSocketConnection(c, db)
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/validate-token", ValidateTokenHandler)

	log.Info().Msg("Starting Server...")
	router.Run()
}

// ValidateTokenHandler validates the JWT token
func ValidateTokenHandler(c *gin.Context) {
	// Get the JWT token from the Authorization header
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization token"})
		return
	}

	// Parse the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("signing-key"), nil // Replace "signing-key" with your actual secret key
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		return
	}

	// Extract user information from the token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
		return
	}

	// Extract any user information you need from the claims
	username, ok := claims["authenticated_user"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username in token"})
		return
	}

	// Return a success response if token is valid
	c.JSON(http.StatusOK, gin.H{"username": username, "message": "Token is valid"})
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
	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal message object to JSON")
		return
	}

	log.Info().Msg(string(messageJSON))

	// Find recipient's WebSocket connection
	recipientConn := userConnections[message.ReceipientID]
	if recipientConn != nil {
		log.Info().Str("recipient_connection_address", recipientConn.RemoteAddr().String()).Msg("WebSocket connection for recipient")
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
		if _, ok := userConnections[receivedMessage.SenderID]; !ok {
			userConnections[receivedMessage.SenderID] = conn
		}

		// Broadcast message to recipient
		sendWebSocketMessage(receivedMessage)
	}
}
