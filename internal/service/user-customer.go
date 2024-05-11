package service

import (
	"context"
	"database/sql"
	"eniqilo-store/internal/domain"
	"eniqilo-store/internal/repository"
	"fmt"
	"reflect"
	"strings"

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
	var whereClause []string
	orderClause := []string{"created_at desc, sid desc"}
	var args []any

	val := reflect.ValueOf(queryParams)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		key := strings.ToLower(typ.Field(i).Name)
		value := val.Field(i).String()
		argPos := len(args) + 1

		if len(value) < 1 {
			continue
		}

		if key == "phonenumber" {
			key = "phone_number"
			value += "%"
		}

		if key == "name" {
			value = "%" + value + "%"
		}

		whereClause = append(whereClause, fmt.Sprintf("%s ILIKE $%d", key, argPos))
		args = append(args, value)
	}

	if len(whereClause) > 0 {
		query += "\nWHERE " + strings.Join(whereClause, " AND ")
	}
	query += "\nORDER BY " + strings.Join(orderClause, ", ")

	customers, err := ucs.userCustomerRepository.GetCustomers(ctx, ucs.db, query, args)
	if err != nil {
		return nil, domain.NewInternalServerError(err.Error())
	}

	return customers, nil
}
