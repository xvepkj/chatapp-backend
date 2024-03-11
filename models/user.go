package models

type User struct {
	UserName string `gorm:"primaryKey;unique"`
	Password string
	Language string
	Token    string `gorm:"-"`

	SentMessages     []Message `gorm:"foreignKey:SenderID"`
	ReceivedMessages []Message `gorm:"foreignKey:ReceipientID"`
}
