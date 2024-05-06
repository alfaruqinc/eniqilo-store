package repository

import "project-sprint-w2/internal/domain"

type UserAdminRepository interface {
	CreateUserAdminRepository() error
	GetUserByIDAdminRepository() (*domain.UserAdmin, error)
	GetUserByPhoneNumberAndPasswordRepository() error
}

type userRepository struct {
}

func NewUserAdminRepository() UserAdminRepository {
	return &userRepository{}
}

func (u *userRepository) CreateUserAdminRepository() error {
	return nil
}

func (u *userRepository) GetUserByIDAdminRepository() (*domain.UserAdmin, error) {
	return &domain.UserAdmin{}, nil
}

func (u *userRepository) GetUserByPhoneNumberAndPasswordRepository() error {
	return nil
}
