package services

import (
	"contract-manage/models"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := models.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := models.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := models.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) GetUsers(skip, limit int) ([]models.User, error) {
	var users []models.User
	if err := models.DB.Offset(skip).Limit(limit).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

type UserCreateInput struct {
	Username   string          `json:"username" binding:"required"`
	Email      string          `json:"email"`
	Password   string          `json:"password" binding:"required"`
	FullName   string          `json:"full_name"`
	Role       models.UserRole `json:"role"`
	Department string          `json:"department"`
	Phone      string          `json:"phone"`
}

func (s *UserService) CreateUser(input UserCreateInput) (*models.User, error) {
	if _, err := s.GetUserByUsername(input.Username); err == nil {
		return nil, errors.New("username already registered")
	}
	if input.Email != "" {
		if _, err := s.GetUserByEmail(input.Email); err == nil {
			return nil, errors.New("email already registered")
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Username:       input.Username,
		Email:          input.Email,
		HashedPassword: string(hashedPassword),
		FullName:       input.FullName,
		Role:           input.Role,
		Department:     input.Department,
		Phone:          input.Phone,
		IsActive:       true,
	}

	if err := models.DB.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

type UserUpdateInput struct {
	Email      string          `json:"email"`
	FullName   string          `json:"full_name"`
	Role       models.UserRole `json:"role"`
	Department string          `json:"department"`
	Phone      string          `json:"phone"`
	IsActive   *bool           `json:"is_active"`
}

func (s *UserService) UpdateUser(id uint, input UserUpdateInput) (*models.User, error) {
	user, err := s.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if input.Email != "" {
		updates["email"] = input.Email
	}
	if input.FullName != "" {
		updates["full_name"] = input.FullName
	}
	if input.Role != "" {
		updates["role"] = input.Role
	}
	if input.Department != "" {
		updates["department"] = input.Department
	}
	if input.Phone != "" {
		updates["phone"] = input.Phone
	}
	if input.IsActive != nil {
		updates["is_active"] = *input.IsActive
	}

	if err := models.DB.Model(user).Updates(updates).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) DeleteUser(id uint) error {
	result := models.DB.Delete(&models.User{}, id)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (s *UserService) AuthenticateUser(username, password string) (*models.User, error) {
	user, err := s.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}
	return user, nil
}