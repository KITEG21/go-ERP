package auth

import "user_api/internal/database"

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) Create(user *User) error {
	return database.DB.Create(user).Error
}

func (r *UserRepository) GetByEmail(email string) (*User, error) {
	var user User
	err := database.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (r *UserRepository) GetByID(id int) (*User, error) {
	var user User
	err := database.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, err
}
