package handlers

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"

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

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}
	user.Password = hashedPassword

	if err := db.CreateUser(dbConn, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := generateJWTToken(user.UserName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	user.Token = token

	c.JSON(http.StatusCreated, gin.H{"message": "user created successfully", "user": user})
}

func GetUser(c *gin.Context, dbConn *gorm.DB) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingUser, err := db.GetUserByUsername(dbConn, user.UserName)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid cridentials"})
		return
	}

	if err := VerifyPassword(existingUser.Password, user.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := generateJWTToken(user.UserName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	existingUser.Token = token
	c.JSON(http.StatusCreated, gin.H{"message": "login successful", "user": existingUser})
}

func generateJWTToken(username string) (string, error) {
	// Create a new JWT token with the appropriate signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"authenticated_user": username,
	})

	// Sign the token with a secret key and get the complete encoded token string
	tokenString, err := token.SignedString([]byte("signing-key"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
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

// HashPassword hashes the given password using bcrypt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// VerifyPassword checks if the provided password matches the hashed password
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
