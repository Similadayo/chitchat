package controller

import (
	"github/similadayo/chitchat/models"
	"github/similadayo/chitchat/utils"
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type MessageController struct {
	DB *gorm.DB
}

// send message
func (mc *MessageController) SendMessage(w http.ResponseWriter, r *http.Request) {
	// Get the sender ID from the context
	senderID, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get the receiver ID from the request
	receiverID, err := strconv.Atoi(r.FormValue("receiver_id"))
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid receiver ID")
		return
	}

	// Get the message content from the request
	content := r.FormValue("content")
	if content == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Message content is required")
		return
	}

	// Convert senderID to uint
	senderIDUint, err := strconv.ParseUint(senderID, 10, 32)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid sender ID")
		return
	}

	// Create a new message
	message := models.Message{
		Content:    content,
		SenderID:   uint(senderIDUint),
		ReceiverID: uint(receiverID),
		Timestamp:  time.Now(),
	}

	// Save the message to the database
	err = mc.DB.Create(&message).Error
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to send message")
		return
	}

	// Respond with the message
	utils.RespondWithJSON(w, http.StatusCreated, message)
}

// get messages
func (mc *MessageController) GetMessages(w http.ResponseWriter, r *http.Request) {
	// Get the sender ID from the context
	senderID, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get the receiver ID from the request
	receiverID, err := strconv.Atoi(r.FormValue("receiver_id"))
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid receiver ID")
		return
	}

	// Get the messages from the database
	var messages []models.Message
	err = mc.DB.Where("sender_id = ? AND receiver_id = ?", senderID, receiverID).Or("sender_id = ? AND receiver_id = ?", receiverID, senderID).Find(&messages).Error
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get messages")
		return
	}

	// Respond with the messages
	utils.RespondWithJSON(w, http.StatusOK, messages)
}
