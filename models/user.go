package models

type User struct {
	UserName string `gorm:"primaryKey;unique"`
	Password string
	Language string

	SentMessages     []Message `gorm:"foreignKey:SenderID"`
	ReceivedMessages []Message `gorm:"foreignKey:ReceipientID"`
}
