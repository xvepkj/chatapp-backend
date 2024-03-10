package models

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	SenderID     string
	ReceipientID string
	Content      string
	Timestamp    time.Time `gorm:"autoCreateTime"`
}
