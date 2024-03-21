package db

import (
	"github.com/xvepkj/chatapp-backend/models"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, user *models.User) error {
	result := db.Create(user)
	return result.Error
}

func GetUserByUsername(db *gorm.DB, username string) (*models.User, error) {
	var user models.User
	result := db.Where("user_name = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func UpdateUser(db *gorm.DB, user *models.User) error {
	result := db.Save(user)
	return result.Error
}

func GetAllUsers(db *gorm.DB) ([]models.User, error) {
	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func DeleteUser(db *gorm.DB, username string) error {
	result := db.Where("user_name = ?", username).Delete(&models.User{})
	return result.Error
}
