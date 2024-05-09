package repository

import (
	"log"

	"github.com/arizdn234/crud-users-and-login-system-with-gin/internal/models"
	"github.com/arizdn234/crud-users-and-login-system-with-gin/utils"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (ur *UserRepository) Create(user *models.User) error {
	return ur.DB.Create(user).Error
}

func (ur *UserRepository) Update(user *models.User) error {
	return ur.DB.Save(user).Error
}

func (ur *UserRepository) Delete(id uint) error {
	return ur.DB.Delete(&models.User{}, id).Error
}

func (ur *UserRepository) FindAll(user *[]models.User) error {
	return ur.DB.Find(user).Error
}

func (ur *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	if err := ur.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := ur.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) Seed() error {
	// Check if there are any users
	var count int64
	if err := ur.DB.Model(&models.User{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	// Create some users
	users := []*models.User{
		{Name: "Johnny Deep", Email: "jdeep557@gmail.com", Password: "blabla2421"},
		{Name: "Tom Cruise", Email: "tcruise557@gmail.com", Password: "bals31dw"},
		{Name: "Billie Elish", Email: "bilelish557@gmail.com", Password: "jup752sdka"},
	}

	for _, user := range users {
		// Hash the password
		hashedPassword, err := utils.HashPassword(user.Password + user.Email)
		if err != nil {
			return err
		}
		user.Password = hashedPassword

		// Create the user
		if err := ur.DB.Create(user).Error; err != nil {
			return err
		}
		log.Println("Seed data successfully added to the database", user)
	}

	return nil
}
