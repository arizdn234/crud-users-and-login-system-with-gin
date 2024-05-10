package handlers

import (
	"net/http"
	"strconv"

	"github.com/arizdn234/crud-users-and-login-system-with-gin/internal/models"
	"github.com/arizdn234/crud-users-and-login-system-with-gin/internal/repository"
	"github.com/arizdn234/crud-users-and-login-system-with-gin/utils"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserRepository *repository.UserRepository
}

func NewUserHandler(ur *repository.UserRepository) *UserHandler {
	return &UserHandler{UserRepository: ur}
}

func (uh *UserHandler) Welcome(c *gin.Context) {
	info := `
	Simple Login & Register system with CRUD on users data!

	Routes Available:
	- GET    /                  : Welcome message
	- GET    /users             : Get all users
	- POST   /users             : Create a new user (using create method)
	- GET    /users/{id}        : Get user by ID
	- PUT    /users/{id}        : Update user by ID
	- DELETE /users/{id}        : Delete user by ID
	- POST   /login             : User login
	- POST   /register          : Register new user (using register method)
	`

	c.String(http.StatusOK, info)
}

func (uh *UserHandler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password + user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user.Password = hashedPassword

	if err := uh.UserRepository.Create(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (uh *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updateReq struct {
		Name     *string `json:"name"`
		Email    *string `json:"email"`
		Password *string `json:"password"`
	}

	if err := c.BindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := uh.UserRepository.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if updateReq.Name != nil {
		user.Name = *updateReq.Name
	}

	if updateReq.Email != nil {
		user.Email = *updateReq.Email
	}

	if updateReq.Password != nil && *updateReq.Password != "" {
		hashedPassword, err := utils.HashPassword(*updateReq.Password + user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		user.Password = hashedPassword
	}

	if err := uh.UserRepository.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (uh *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := uh.UserRepository.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
