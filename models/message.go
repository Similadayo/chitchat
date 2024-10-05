package models

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	Content     string    `json:"content" gorm:"not null"`
	ImageURL    string    `json:"image_url"`
	IsDelivered bool      `json:"is_delivered" gorm:"default:false"`
	IsRead      bool      `json:"is_read" gorm:"default:false"`
	IsEdited    bool      `json:"is_edited" gorm:"default:false"`
	IsDeleted   bool      `json:"is_deleted" gorm:"default:false"`
	SenderID    uint      `json:"sender_id" gorm:"not null"`
	ReceiverID  uint      `json:"receiver_id" gorm:"not null"`
	GroupID     uint      `json:"group_id"`
	Timestamp   time.Time `json:"timestamp" gorm:"not null"`
}
