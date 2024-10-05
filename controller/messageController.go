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
		Content:     content,
		ImageURL:    r.FormValue("image_url"),
		IsDelivered: false, // Initially set to false when sent
		IsRead:      false, // Initially set to false when sent
		IsEdited:    false,
		IsDeleted:   false,
		SenderID:    uint(senderIDUint),
		ReceiverID:  uint(receiverID),
		Timestamp:   time.Now(),
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
	err = mc.DB.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", senderID, receiverID, receiverID, senderID).Find(&messages).Error
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get messages")
		return
	}

	// Respond with the messages
	utils.RespondWithJSON(w, http.StatusOK, messages)
}

// MarkAsRead marks a message as read
func (mc *MessageController) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	// Get the receiver ID from the context
	receiverID, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get the message ID from the request
	messageID, err := strconv.Atoi(r.FormValue("message_id"))
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid message ID")
		return
	}

	// Get the message from the database
	var message models.Message
	err = mc.DB.Where("id = ? AND receiver_id = ?", messageID, receiverID).First(&message).Error
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Message not found")
		return
	}

	// Mark the message as read
	message.IsRead = true

	// Save the message to the database
	err = mc.DB.Save(&message).Error
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to mark message as read")
		return
	}

	// Respond with the message
	utils.RespondWithJSON(w, http.StatusOK, message)
}

// MarkAsDelivered marks a message as delivered
func (mc *MessageController) MarkAsDelivered(w http.ResponseWriter, r *http.Request) {
	// Get the receiver ID from the context
	receiverID, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get the message ID from the request
	messageID, err := strconv.Atoi(r.FormValue("message_id"))
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid message ID")
		return
	}

	// Get the message from the database
	var message models.Message
	err = mc.DB.Where("id = ? AND receiver_id = ?", messageID, receiverID).First(&message).Error
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Message not found")
		return
	}

	// Mark the message as delivered
	message.IsDelivered = true

	// Save the message to the database
	err = mc.DB.Save(&message).Error
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to mark message as delivered")
		return
	}

	// Respond with the message
	utils.RespondWithJSON(w, http.StatusOK, message)
}

// EditMessage edits a message
func (mc *MessageController) EditMessage(w http.ResponseWriter, r *http.Request) {
	// Get the sender ID from the context
	senderID, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get the message ID from the request
	messageID, err := strconv.Atoi(r.FormValue("message_id"))
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid message ID")
		return
	}

	// Get the new message content from the request
	content := r.FormValue("content")
	if content == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Message content is required")
		return
	}

	// Get the message from the database
	var message models.Message
	err = mc.DB.Where("id = ? AND sender_id = ?", messageID, senderID).First(&message).Error
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Message not found")
		return
	}

	// Update the message content
	message.Content = content

	// Save the message to the database
	err = mc.DB.Save(&message).Error
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to edit message")
		return
	}

	// Respond with the message
	utils.RespondWithJSON(w, http.StatusOK, message)
}

// DeleteMessage deletes a message
func (mc *MessageController) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	// Get the sender ID from the context
	senderID, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get the message ID from the request
	messageID, err := strconv.Atoi(r.FormValue("message_id"))
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid message ID")
		return
	}

	// Get the message from the database
	var message models.Message
	err = mc.DB.Where("id = ? AND sender_id = ?", messageID, senderID).First(&message).Error
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Message not found")
		return
	}

	// Delete the message from the database
	err = mc.DB.Delete(&message).Error
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete message")
		return
	}

	// Respond with success
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Message deleted"})
}

// SearchMessages searches for messages
func (mc *MessageController) SearchMessages(w http.ResponseWriter, r *http.Request) {
	// Get the sender ID from the context
	senderID, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get the search query from the request
	query := r.FormValue("query")
	if query == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Search query is required")
		return
	}

	// Get the messages from the database
	var messages []models.Message
	err := mc.DB.Where("(sender_id = ? OR receiver_id = ?) AND content LIKE ?", senderID, senderID, "%"+query+"%").Find(&messages).Error
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to search messages")
		return
	}

	// Respond with the messages
	utils.RespondWithJSON(w, http.StatusOK, messages)
}
