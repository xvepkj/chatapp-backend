package db

import (
	"github.com/xvepkj/chatapp-backend/models"
	"gorm.io/gorm"
)

func AddMessage(db *gorm.DB, message *models.Message) error {
	result := db.Create(message)
	return result.Error
}

func GetMessagesBetween(db *gorm.DB, senderID string, receiverID string) ([]models.Message, error) {
	var messages []models.Message

	result := db.Where("(sender_id = ? AND receipient_id = ?) OR (sender_id = ? AND receipient_id = ?)",
		senderID, receiverID, receiverID, senderID).Find(&messages)

	if result.Error != nil {
		return nil, result.Error
	}

	return messages, nil
}
