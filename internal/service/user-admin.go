package service

import (
	"context"
	"database/sql"
	"project-sprint-w2/internal/domain"
	"project-sprint-w2/internal/repository"
)

type UserAdminService interface {
	RegisterUserAdminService(ctx context.Context, userAdmin domain.RegisterUserAdmin) (*domain.UserAdminResponseWithAccessToken, error)
	LoginUserAdminService(ctx context.Context, userAdmin domain.LoginUserAdmin) (*domain.UserAdminResponseWithAccessToken, error)
}

type userAdminService struct {
	db                  *sql.DB
	userAdminRepository repository.UserAdminRepository
}

func NewUserAdminService(db *sql.DB, userAdminRepository repository.UserAdminRepository) UserAdminService {
	return &userAdminService{
		db:                  db,
		userAdminRepository: userAdminRepository,
	}
}

func (u *userAdminService) RegisterUserAdminService(ctx context.Context, userAdmin domain.RegisterUserAdmin) (*domain.UserAdminResponseWithAccessToken, error) {
	return &domain.UserAdminResponseWithAccessToken{}, nil
}

func (u *userAdminService) LoginUserAdminService(ctx context.Context, userAdmin domain.LoginUserAdmin) (*domain.UserAdminResponseWithAccessToken, error) {
	return &domain.UserAdminResponseWithAccessToken{}, nil
}
