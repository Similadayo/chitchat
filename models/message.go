package models

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	Content    string    `json:"content" gorm:"not null"`
	SenderID   uint      `json:"sender_id" gorm:"not null"`
	ReceiverID uint      `json:"receiver_id" gorm:"not null"`
	GroupID    uint      `json:"group_id"`
	Timestamp  time.Time `json:"timestamp" gorm:"not null"`
}
