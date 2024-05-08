package repository

import (
	"context"
	"database/sql"
	"eniqilo-store/internal/domain"
)

type UserAdminRepository interface {
	CreateUserAdminRepository(ctx context.Context, db *sql.DB, userAdmin domain.UserAdmin) error
	GetUserByIDAdminRepository(ctx context.Context, tx *sql.Tx, id string) (*domain.UserAdmin, error)
	GetUserByPhoneNumberRepository(ctx context.Context, db *sql.DB, phoneNumber string) (*domain.UserAdmin, error)
	CheckPhoneNumberExists(ctx context.Context, tx *sql.Tx, phoneNumber string) (bool, error)
}

type userRepository struct{}

func NewUserAdminRepository() UserAdminRepository {
	return &userRepository{}
}

func (u *userRepository) CreateUserAdminRepository(ctx context.Context, db *sql.DB, userAdmin domain.UserAdmin) error {
	query := `INSERT INTO user_admins (id, created_at, phone_number, password, name, role) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := db.ExecContext(ctx, query, userAdmin.ID, userAdmin.CreatedAt, userAdmin.PhoneNumber, userAdmin.Password, userAdmin.Name, userAdmin.Role)
	if err != nil {
		return err
	}

	return nil
}

func (u *userRepository) GetUserByIDAdminRepository(ctx context.Context, tx *sql.Tx, id string) (*domain.UserAdmin, error) {
	user := domain.UserAdmin{}

	query := `SELECT id, phone_number, name, role FROM user_admins WHERE id = $1`

	row := tx.QueryRowContext(ctx, query, id)
	err := row.Scan(user)
	if err != nil {
		return nil, err
	}

	return &domain.UserAdmin{}, nil
}

func (u *userRepository) GetUserByPhoneNumberRepository(ctx context.Context, db *sql.DB, phoneNumber string) (*domain.UserAdmin, error) {
	user := domain.UserAdmin{}

	query := `SELECT id, created_at, phone_number, name, role, password FROM user_admins WHERE phone_number = $1`

	row := db.QueryRowContext(ctx, query, phoneNumber)
	err := row.Scan(&user.ID, &user.CreatedAt, &user.PhoneNumber, &user.Name, &user.Role, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *userRepository) CheckPhoneNumberExists(ctx context.Context, tx *sql.Tx, phoneNumber string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM user_admins
			WHERE phone_number = $1
		)
	`
	var exists bool
	err := tx.QueryRowContext(ctx, query, phoneNumber).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
