package routes

import (
	"github/similadayo/chitchat/controller"

	"github.com/gorilla/mux"
)

func InitRoutes() *mux.Router {
	router := mux.NewRouter()

	// Creates new user controller
	userController := controller.NewUserController()

	// Add routes here
	router.HandleFunc("/register", userController.RegisterUser).Methods("POST")
	router.HandleFunc("/login", userController.LoginUser).Methods("POST")

	// Add routes for the user profile

	return router

}
