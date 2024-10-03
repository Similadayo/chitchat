package controller

import (
	"encoding/json"
	"github/similadayo/chitchat/config"
	"github/similadayo/chitchat/models"
	"github/similadayo/chitchat/utils"
	"net/http"

	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

// NewUserController creates a new user controller
func NewUserController() *UserController {
	db, err := config.ConnectDB()
	if err != nil {
		panic(err)
	}

	config.MigrateDB(db)
	return &UserController{DB: db}
}

// RegisterUser registers a new user
func (uc *UserController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	//Parse the request body
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	//Hash the password before saving the user
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Could not hash the password", http.StatusInternalServerError)
		return
	}
	user.Password = hashedPassword

	//handles image upload for profile picture
	file, handler, err := r.FormFile("profile_pic")
	if err == nil {
		defer file.Close()
		user.ProfilePic = handler.Filename
	}

	//Save the user in the database
	if err := uc.DB.Create(&user).Error; err != nil {
		http.Error(w, "Could not save the user", http.StatusInternalServerError)
		return
	}

	//Respond with success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}

// LoginUser logs in a user
func (uc *UserController) LoginUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	//Parse the request body
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	//Check if the user exists
	var existingUser models.User
	if err := uc.DB.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
		http.Error(w, "Invalid credentials email", http.StatusUnauthorized)
		return
	}

	//Check if the password is correct
	if err := utils.ComparePasswords(existingUser.Password, user.Password); err != nil {
		http.Error(w, "Invalid credentials password", http.StatusUnauthorized)
		return
	}

	//Generate a JWT token for the authenticated User
	token, err := utils.GenerateJwt(existingUser.Username)
	if err != nil {
		http.Error(w, "Could not generate JWT token", http.StatusInternalServerError)
		return
	}

	//Respond with the token
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
