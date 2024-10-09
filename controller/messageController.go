package controller

import (
	"encoding/json"
	"github/similadayo/chitchat/config"
	"github/similadayo/chitchat/models"
	"github/similadayo/chitchat/utils"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type MessageController struct {
	DB *gorm.DB
}

func NewMessageController() *MessageController {
	db, err := config.ConnectDB()
	if err != nil {
		panic(err)
	}

	config.MigrateDB(db)

	return &MessageController{DB: db}
}

// Func SendMessage is used to send a message to a user or a group
func (mc *MessageController) SendMessage(w http.ResponseWriter, r *http.Request) {
	var message models.Message
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	senderID, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var sender models.User
	if err := mc.DB.Where("username = ?", senderID).First(&sender).Error; err != nil {
		http.Error(w, "Sender not found", http.StatusNotFound)
		return
	}

	message.SenderID = sender.ID
	message.Timestamp = time.Now()

	//validate receiver
	var receiver models.User
	if err := mc.DB.First(&receiver, message.ReceiverID).Error; err != nil {
		http.Error(w, "Receiver not found", http.StatusNotFound)
		return
	}

	if err := mc.DB.Create(&message).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Message sent successfully"})
}

// Func GetMessages is used to get all messages sent to a user or a group
func (mc *MessageController) GetMessages(w http.ResponseWriter, r *http.Request) {
	senderID := r.URL.Query().Get("sender_id")
	receiverID := r.URL.Query().Get("receiver_id")

	receiverIDuint, err := utils.ConvertToUint(receiverID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var messages []models.Message
	if err := mc.DB.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", senderID, receiverID, receiverID, senderID).
		Preload("Sender").
		Preload("Receiver").
		Find(&messages).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update the IsDelivered field for messages sent to the current user
	for i := range messages {
		if messages[i].ReceiverID == uint(receiverIDuint) && !messages[i].IsDelivered {
			messages[i].IsDelivered = true
			if err := mc.DB.Model(&messages[i]).Update("is_delivered", true).Error; err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
