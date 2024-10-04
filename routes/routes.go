package routes

import (
	"github/similadayo/chitchat/controller"
	"github/similadayo/chitchat/middlewares"

	"github.com/gorilla/mux"
)

func InitRoutes() *mux.Router {
	router := mux.NewRouter()

	// Creates new user controller
	userController := controller.NewUserController()

	// Add routes here
	router.HandleFunc("/register", userController.RegisterUser).Methods("POST")
	router.HandleFunc("/login", userController.LoginUser).Methods("POST")

	// protected routes
	protected := router.PathPrefix("/").Subrouter()
	protected.Use(middlewares.AuthMiddleware)
	protected.HandleFunc("/users", userController.GetAllUsers).Methods("GET")
	protected.HandleFunc("/user", userController.GetUserProfile).Methods("GET")
	protected.HandleFunc("/users/{username}", userController.GetUserByUserName).Methods("GET")
	protected.HandleFunc("/user/update", userController.UpdateUserProfile).Methods("PUT")
	protected.HandleFunc("/user/delete", userController.DeleteUserProfile).Methods("DELETE")
	protected.HandleFunc("/block/{username}", userController.BlockUser).Methods("POST")
	protected.HandleFunc("/unblock/{username}", userController.UnblockUser).Methods("POST")
	protected.HandleFunc("/logout", userController.Logout).Methods("POST")

	return router

}
