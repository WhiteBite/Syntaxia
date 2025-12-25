package services

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	ID    uint   `json:"id" gorm:"primary_key"`
	Name  string `json:"name"`
	Email string `json:"email" gorm:"unique"`
}

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) GetAllUsers() ([]User, error) {
	var users []User
	if err := s.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) GetUserByID(id uint) (*User, error) {
	var user User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) CreateUser(user *User) error {
	return s.db.Create(user).Error
}

func (s *UserService) UpdateUser(user *User) error {
	return s.db.Save(user).Error
}

func (s *UserService) DeleteUser(id uint) error {
	return s.db.Delete(&User{}, id).Error
}
