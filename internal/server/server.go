package server

import (
	"github.com/arizdn234/crud-users-and-login-system-with-gin/internal/handlers"
	"github.com/arizdn234/crud-users-and-login-system-with-gin/internal/middleware"
	"github.com/arizdn234/crud-users-and-login-system-with-gin/internal/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RunServer(r *gin.Engine, db *gorm.DB, port string) *gin.Engine {
	userRepo := repository.NewUserRepository(db)
	userHandler := handlers.NewUserHandler(userRepo)

	userRoute := r.Group("/users")
	{

		userRoute.GET("/", userHandler.Welcome)
		userRoute.POST("/users/login", userHandler.UserLogin)
		userRoute.POST("/users/register", userHandler.UserRegister)
		userRoute.GET("/users/logout", userHandler.UserLogout)

		userRoute.Use(middleware.RequireAuth()) // endpoints that require tokens from this endpoint group
		userRoute.GET("", userHandler.GetAllUsers)
		userRoute.GET("/:id", userHandler.GetUserByID)
		userRoute.POST("", userHandler.CreateUser)
		userRoute.PUT("/:id", userHandler.UpdateUser)
		userRoute.DELETE("/:id", userHandler.DeleteUser)
	}

	return r
}
