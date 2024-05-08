package repository

import (
	"context"
	"database/sql"
	"eniqilo-store/internal/domain"
)

type UserCustomerRepository interface {
	CreateUserCustomer(ctx context.Context, db *sql.DB, userCustomer domain.UserCustomer) error
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
