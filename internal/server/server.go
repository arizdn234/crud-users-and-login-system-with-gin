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

	r.GET("/", userHandler.Welcome)

	userRoute := r.Group("/users")
	{

		userRoute.POST("/login", userHandler.UserLogin)
		userRoute.POST("/register", userHandler.UserRegister)
		userRoute.GET("/logout", userHandler.UserLogout)

		userRoute.Use(middleware.RequireAuth()) // endpoints that require tokens from this endpoint group
		userRoute.GET("", userHandler.GetAllUsers)
		userRoute.GET("/:id", userHandler.GetUserByID)
		userRoute.POST("", userHandler.CreateUser)
		userRoute.PUT("/:id", userHandler.UpdateUser)
		userRoute.DELETE("/:id", userHandler.DeleteUser)
	}

	return r
}
