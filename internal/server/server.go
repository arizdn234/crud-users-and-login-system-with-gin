package server

import (
	"net/http"

	"github.com/arizdn234/crud-users-and-login-system-with-gin/internal/middleware"
	"github.com/arizdn234/crud-users-and-login-system-with-gin/internal/repository"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func NewServer(db *gorm.DB, port string) *http.Server {
	userRepo := repository.NewUserRepository(db)
	userHandler := handlers.NewUserHandler(userRepo)

	r := mux.NewRouter()

	// no need login
	r.HandleFunc("/", userHandler.Welcome).Methods(http.MethodGet)
	r.HandleFunc("/users/login", userHandler.UserLogin).Methods(http.MethodPost)
	r.HandleFunc("/users/register", userHandler.UserRegister).Methods(http.MethodPost)
	r.HandleFunc("/users/logout", userHandler.UserLogout).Methods(http.MethodGet)

	// need login
	authorizedRoute := r.PathPrefix("/users").Subrouter()
	authorizedRoute.Use(middleware.RequireAuth)
	authorizedRoute.HandleFunc("", userHandler.GetAllUsers).Methods(http.MethodGet)
	authorizedRoute.HandleFunc("/{id}", userHandler.GetUserByID).Methods(http.MethodGet)
	authorizedRoute.HandleFunc("", userHandler.CreateUser).Methods(http.MethodPost)
	authorizedRoute.HandleFunc("/{id}", userHandler.UpdateUser).Methods(http.MethodPut)
	authorizedRoute.HandleFunc("/{id}", userHandler.DeleteUser).Methods(http.MethodDelete)

	return &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}
}
