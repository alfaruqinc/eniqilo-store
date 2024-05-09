package repository

import (
	"context"
	"database/sql"
	"eniqilo-store/internal/domain"
	"fmt"
	"reflect"
	"strings"
)

type CheckoutRepository interface {
	CreateCheckout(ctx context.Context, tx *sql.Tx, checkout domain.Checkout, productCheckout []domain.ProductCheckout) error
	GetCheckoutHistory(ctx context.Context, db *sql.DB) ([]domain.GetCheckoutHistory, error)
}

type checkoutRepository struct{}

func NewCheckoutRepository() CheckoutRepository {
	return &checkoutRepository{}
}

func (cr *checkoutRepository) CreateCheckout(ctx context.Context, tx *sql.Tx, checkout domain.Checkout, productCheckouts []domain.ProductCheckout) error {
	query := `
		INSERT INTO checkouts (id, created_at, user_customer_id, paid, change)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := tx.ExecContext(ctx, query, checkout.ID, checkout.CreatedAt, checkout.UserCustomerID, checkout.Paid, checkout.Change)
	if err != nil {
		return err
	}

	inserts := []string{}
	args := []any{}
	for _, pc := range productCheckouts {
		val := reflect.ValueOf(pc)

		insert := []string{}
		for i := 0; i < val.NumField(); i++ {
			value := val.Field(i).Interface()
			argsPos := len(args) + 1

			insert = append(insert, fmt.Sprintf("$%d", argsPos))
			args = append(args, value)
		}
		placeholder := fmt.Sprintf("(%s)", strings.Join(insert, ", "))
		inserts = append(inserts, placeholder)
	}
	query = `
		INSERT INTO product_checkouts (id, product_id, quantity, checkout_id)
		VALUES `
	query += strings.Join(inserts, ", ")

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (cr *checkoutRepository) GetCheckoutHistory(ctx context.Context, db *sql.DB) ([]domain.GetCheckoutHistory, error) {
	query := `
		SELECT c.id, c.user_customer_id, pc.product_id, pc.quantity, c.paid, c.change
		FROM checkouts c
		INNER JOIN product_checkouts pc ON pc.checkout_id = c.id
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	checkouts := []domain.GetCheckoutHistory{}
	for rows.Next() {
		checkout := domain.GetCheckoutHistory{}

		err := rows.Scan(&checkout.TransactionID, &checkout.CustomerID, &checkout.ProductID, &checkout.Quantity, &checkout.Paid, &checkout.Change)
		if err != nil {
			return nil, err
		}

		checkouts = append(checkouts, checkout)
	}

	return checkouts, nil
}
