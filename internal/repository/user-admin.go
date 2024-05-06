package repository

import (
	"context"
	"database/sql"
	"project-sprint-w2/internal/domain"
)

type UserAdminRepository interface {
	CreateUserAdminRepository(ctx context.Context, tx *sql.Tx, userAdmin domain.UserAdmin) error
	GetUserByIDAdminRepository(ctx context.Context, tx *sql.Tx, id string) (*domain.UserAdmin, error)
	GetUserByPhoneNumberRepository(ctx context.Context, tx *sql.Tx, email string) (*domain.UserAdmin, error)
}

type userRepository struct {
}

func NewUserAdminRepository() UserAdminRepository {
	return &userRepository{}
}

func (u *userRepository) CreateUserAdminRepository(ctx context.Context, tx *sql.Tx, userAdmin domain.UserAdmin) error {
	query := `INSERT INTO user_admin (phone_number, password, name, role) VALUES ($1, $2, $3, $4)`

	_, err := tx.ExecContext(ctx, query, userAdmin.PhoneNumber, userAdmin.Password, userAdmin.Name, userAdmin.Role)
	if err != nil {
		return err
	}

	return nil
}

func (u *userRepository) GetUserByIDAdminRepository(ctx context.Context, tx *sql.Tx, id string) (*domain.UserAdmin, error) {
	user := domain.UserAdmin{}

	query := `SELECT id, phone_number, name, role FROM user_admin WHERE id = $1`

	row := tx.QueryRowContext(ctx, query, id)
	err := row.Scan(user)
	if err != nil {
		return nil, err
	}

	return &domain.UserAdmin{}, nil
}

func (u *userRepository) GetUserByPhoneNumberRepository(ctx context.Context, tx *sql.Tx, phoneNumber string) (*domain.UserAdmin, error) {
	user := domain.UserAdmin{}

	query := `SELECT id, phone_number, name, role, password FROM user_admin WHERE phone_number = $1`

	row := tx.QueryRowContext(ctx, query, phoneNumber)
	err := row.Scan(user)
	if err != nil {
		return nil, err
	}

	return &domain.UserAdmin{}, nil
}
