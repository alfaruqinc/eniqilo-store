package service

import (
	"context"
	"database/sql"
	"eniqilo-store/internal/domain"
	"eniqilo-store/internal/repository"
)

type UserCustomerService interface {
	CreateUserCustomer(ctx context.Context, userCustomer domain.UserCustomer) domain.MessageErr
}

type userCustomerService struct {
	db                     *sql.DB
	userCustomerRepository repository.UserCustomerRepository
}

func NewUserCustomerService(db *sql.DB, userCustomerRepository repository.UserCustomerRepository) UserCustomerService {
	return &userCustomerService{
		db:                     db,
		userCustomerRepository: userCustomerRepository,
	}
}

func (ucs *userCustomerService) CreateUserCustomer(ctx context.Context, userCustomer domain.UserCustomer) domain.MessageErr {
	err := ucs.userCustomerRepository.CreateUserCustomer(ctx, ucs.db, userCustomer)
	if err != nil {
		return domain.NewInternalServerError(err.Error())
	}

	return nil
}
