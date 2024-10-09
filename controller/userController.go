package controller

import (
	"encoding/json"
	"fmt"
	"github/similadayo/chitchat/config"
	"github/similadayo/chitchat/models"
	"github/similadayo/chitchat/utils"
	"net/http"
	"time"

	"github.com/gorilla/mux"
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

	// Parse the request body
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Hash the password before saving the user
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Could not hash the password", http.StatusInternalServerError)
		return
	}
	user.Password = hashedPassword

	// Handle image upload for profile picture
	file, handler, err := r.FormFile("profile_pic")
	if err == nil {
		defer file.Close()
		user.ProfilePic = handler.Filename
	}

	// Save the user in the database
	if err := uc.DB.Create(&user).Error; err != nil {
		http.Error(w, "Could not save the user", http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}

// LoginUser logs in a user
func (uc *UserController) LoginUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	// Parse the request body
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if the user exists
	var existingUser models.User
	if err := uc.DB.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Check if the password is correct
	if err := utils.ComparePasswords(existingUser.Password, user.Password); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate a JWT token for the authenticated user
	token, err := utils.GenerateJwt(existingUser.Username)
	if err != nil {
		http.Error(w, "Could not generate JWT token", http.StatusInternalServerError)
		return
	}

	// Respond with the token
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// Other functions (e.g., update user, delete user, etc.) can be added similarly.

// GetUserProfile fetches the profile of the currently logged-in user
func (uc *UserController) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	// Retrieve the username from the request context
	username, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		fmt.Println("Error: Could not extract user from context")
		http.Error(w, "Could not extract user from context", http.StatusInternalServerError)
		return
	}
	fmt.Println("Username from context:", username)

	// Fetch the full user from the database using the username
	var user models.User
	if err := uc.DB.Where("username = ?", username).First(&user).Error; err != nil {
		fmt.Println("Error fetching user from database:", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Respond with the user profile
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (uc *UserController) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	username, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		fmt.Println("Error: Could not extract user from context")
		http.Error(w, "Could not extract user from context", http.StatusInternalServerError)
		return
	}

	var user models.User
	if err := uc.DB.Where("username = ?", username).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var updatedUser models.User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	//handles image upload for profile picture
	file, handler, err := r.FormFile("profile_pic")
	if err == nil {
		defer file.Close()
		updatedUser.ProfilePic = handler.Filename
	}

	user.FirstName = updatedUser.FirstName
	user.LastName = updatedUser.LastName
	user.Email = updatedUser.Email
	user.ProfilePic = updatedUser.ProfilePic

	if err := uc.DB.Save(&user).Error; err != nil {
		http.Error(w, "Could not update the user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (uc *UserController) DeleteUserProfile(w http.ResponseWriter, r *http.Request) {
	username, ok := utils.GetUserFromContext(r.Context())
	if !ok {
		fmt.Println("Error: Could not extract user from context")
		http.Error(w, "Could not extract user from context", http.StatusInternalServerError)
		return
	}

	var user models.User
	if err := uc.DB.Where("username = ?", username).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if err := uc.DB.Delete(&user).Error; err != nil {
		http.Error(w, "Could not delete the user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
}

// GetAllUsers returns all the users
func (uc *UserController) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User

	//Get all the users from the database
	if err := uc.DB.Find(&users).Error; err != nil {
		http.Error(w, "Could not get the users", http.StatusInternalServerError)
		return
	}

	//Respond with the users
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

// GetUserByUserName returns a user by UserName
func (uc *UserController) GetUserByUserName(w http.ResponseWriter, r *http.Request) {
	//Get the UserName from the URL
	vars := mux.Vars(r)
	userName := vars["username"]

	//Check if the UserName is empty
	if userName == "" {
		http.Error(w, "UserName is required", http.StatusBadRequest)
		return
	}

	//Get the user by UserName
	var user models.User
	if err := uc.DB.Where("username = ?", userName).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	//Respond with the user
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// GetUserByEmail returns a user by email
func (uc *UserController) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	//Get the email from the URL
	email := r.URL.Query().Get("email")

	//Check if the email is empty
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	//Get the user by email
	var user models.User
	if err := uc.DB.Where("email = ?", email).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	//Respond with the user
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// BlockUser blocks a user
func (uc *UserController) BlockUser(w http.ResponseWriter, r *http.Request) {
	//Get the authenticated user from the request context
	user := r.Context().Value("user").(models.User)

	//Get the user to block by UserName
	userName := r.URL.Query().Get("username")

	//Check if the UserName is empty
	if userName == "" {
		http.Error(w, "UserName is required", http.StatusBadRequest)
		return
	}

	//Get the user to block
	var userToBlock models.User
	if err := uc.DB.Where("username = ?", userName).First(&userToBlock).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	//Block the user
	if err := uc.DB.Model(&user).Association("BlockedUsers").Append(&userToBlock); err != nil {
		http.Error(w, "Could not block the user", http.StatusInternalServerError)
		return
	}

	//Respond with success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User blocked successfully"})
}

// UnblockUser unblocks a user
func (uc *UserController) UnblockUser(w http.ResponseWriter, r *http.Request) {
	//Get the authenticated user from the request context
	user := r.Context().Value("user").(models.User)

	//Get the user to unblock by UserName
	userName := r.URL.Query().Get("username")

	//Check if the UserName is empty
	if userName == "" {
		http.Error(w, "UserName is required", http.StatusBadRequest)
		return
	}

	//Get the user to unblock
	var userToUnblock models.User
	if err := uc.DB.Where("username = ?", userName).First(&userToUnblock).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	//Unblock the user
	if err := uc.DB.Model(&user).Association("BlockedUsers").Delete(&userToUnblock); err != nil {
		http.Error(w, "Could not unblock the user", http.StatusInternalServerError)
		return
	}

	//Respond with success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User unblocked successfully"})
}

// LogoutUser logs out a user
func (uc *UserController) Logout(w http.ResponseWriter, r *http.Request) {
	// Clear the Authorization header or token on client side
	w.Header().Set("Authorization", "")

	// You could also delete a token from cookies (if stored in cookies)
	http.SetCookie(w, &http.Cookie{
		Name:     "Authorization",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour), // Expiring the cookie immediately
		HttpOnly: true,
	})

	// Optionally respond with a message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully logged out"))
}
