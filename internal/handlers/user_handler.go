package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

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

	// validate user data
	if err := uh.validateUser(user); err != nil {
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

	var updatePayload struct {
		Name     *string `json:"name"`
		Email    *string `json:"email"`
		Password *string `json:"password"`
	}

	if err := c.BindJSON(&updatePayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userUpdate = models.User{
		Name:     *updatePayload.Name,
		Email:    *updatePayload.Email,
		Password: *updatePayload.Password,
	}

	// validate user data
	if err := uh.validateUser(userUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := uh.UserRepository.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if updatePayload.Name != nil {
		user.Name = *updatePayload.Name
	}

	if updatePayload.Email != nil {
		user.Email = *updatePayload.Email
	}

	if updatePayload.Password != nil && *updatePayload.Password != "" {
		hashedPassword, err := utils.HashPassword(*updatePayload.Password + user.Email)
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

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Delete user with ID=%v success", id)})
}

func (uh *UserHandler) GetAllUsers(c *gin.Context) {
	var users []models.User

	if err := uh.UserRepository.FindAll(&users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (uh *UserHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := uh.UserRepository.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (uh *UserHandler) UserRegister(c *gin.Context) {
	var registerPayload models.UserRegister
	if err := c.BindJSON(&registerPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user = models.User{
		Name:     registerPayload.Name,
		Email:    registerPayload.Email,
		Password: registerPayload.Password,
	}

	// validate user data
	if err := uh.validateUser(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// email check
	existingUser, err := uh.UserRepository.FindByEmail(user.Email)
	if err == nil && existingUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already registered"})
		return
	}

	// hash
	hashed, err := utils.HashPassword(user.Password + user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user.Password = hashed

	if err := uh.UserRepository.Create(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (uh *UserHandler) UserLogin(c *gin.Context) {
	var loginPayload models.UserLogin
	if err := c.BindJSON(&loginPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userLogin = models.User{
		Email:    loginPayload.Email,
		Password: loginPayload.Password,
	}

	// validate user data
	if err := uh.validateUser(userLogin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if loginPayload.Email == "" || loginPayload.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Login data is empty"})
		return
	}

	existingUser, err := uh.UserRepository.FindByEmail(loginPayload.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	hashedRequestBodyPassword, err := utils.HashPassword(loginPayload.Password + loginPayload.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(hashedRequestBodyPassword)
	fmt.Printf("existingUser.Password: %v\n", existingUser.Password)

	if hashedRequestBodyPassword != existingUser.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// user credentials have been verified
	// generate jwt token
	token, err := utils.CreateToken(&existingUser.ID, &existingUser.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("auth_token", *token, int(time.Hour.Seconds()*24), "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "login successful"})
}

func (uh *UserHandler) UserLogout(c *gin.Context) {
	c.SetCookie("auth_token", "", -1, "/", "", false, true)
	c.String(http.StatusOK, "logout successful")
}

func (uh *UserHandler) validateUser(user models.User) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(user.Email) {
		return errors.New("invalid email format")
	}

	if len(user.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString(user.Password)
	hasLowercase := regexp.MustCompile(`[a-z]`).MatchString(user.Password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(user.Password)

	if !hasUppercase || !hasLowercase || !hasNumber {
		return errors.New("password must contain at least one uppercase letter, one lowercase letter, and one number")
	}

	return nil
}
