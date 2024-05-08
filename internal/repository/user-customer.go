package repository

import (
	"context"
	"database/sql"
	"eniqilo-store/internal/domain"
)

type UserCustomerRepository interface {
	CreateUserCustomer(ctx context.Context, db *sql.DB, userCustomer domain.UserCustomer) error
	GetCustomers(ctx context.Context, db *sql.DB, queryParams string, args []any) ([]domain.UserCustomerResponse, error)
}

type userCustomerRepository struct{}

func NewUserCustomerRepository() UserCustomerRepository {
	return &userCustomerRepository{}
}

func (ucr *userCustomerRepository) CreateUserCustomer(ctx context.Context, db *sql.DB, userCustomer domain.UserCustomer) error {
	query := `
		INSERT INTO user_customers (id, created_at, name, phone_number)
		VALUES ($1, $2, $3, $4)
	`
	_, err := db.ExecContext(ctx, query, userCustomer.ID, userCustomer.CreatedAt, userCustomer.Name, userCustomer.PhoneNumber)
	if err != nil {
		return err
	}

	return nil
}

func (ucr *userCustomerRepository) GetCustomers(ctx context.Context, db *sql.DB, queryParams string, args []any) ([]domain.UserCustomerResponse, error) {
	query := `
		SELECT id, phone_number, name
		FROM user_customers
	`
	query += queryParams

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []domain.UserCustomerResponse
	for rows.Next() {
		customer := domain.UserCustomerResponse{}

		err := rows.Scan(&customer.ID, &customer.PhoneNumber, &customer.Name)
		if err != nil {
			return nil, err
		}

		customers = append(customers, customer)
	}

	return customers, nil
}
