package service

import (
	"context"
	"database/sql"
	"eniqilo-store/internal/domain"
	"eniqilo-store/internal/repository"

	"github.com/jackc/pgx/v5/pgconn"
)

type UserCustomerService interface {
	CreateUserCustomer(ctx context.Context, userCustomer domain.UserCustomer) domain.MessageErr
	GetUserCustomers(ctx context.Context, queryParams domain.UserCustomerQueryParams) ([]domain.UserCustomerResponse, domain.MessageErr)
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
		if err, ok := err.(*pgconn.PgError); ok {
			if err.Code == "23505" {
				return domain.NewConflictError("phone number already exists")
			}
		}

		return domain.NewInternalServerError(err.Error())
	}

	return nil
}

func (ucs *userCustomerService) GetUserCustomers(ctx context.Context, queryParams domain.UserCustomerQueryParams) ([]domain.UserCustomerResponse, domain.MessageErr) {
	var query string
	var args []any

	customers, err := ucs.userCustomerRepository.GetCustomers(ctx, ucs.db, query, args)
	if err != nil {
		return nil, domain.NewInternalServerError(err.Error())
	}

	return customers, nil
}
