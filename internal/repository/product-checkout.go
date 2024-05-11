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
	CreateCheckout(ctx context.Context, tx *sql.Tx, checkout domain.Checkout, productCheckouts []domain.ProductCheckout) error
	GetCheckoutHistory(ctx context.Context, db *sql.DB, queryParams domain.CheckoutHistoryQueryParams) ([]domain.GetCheckoutHistory, error)
	CreateProductCheckout(ctx context.Context, tx *sql.Tx, productCheckout domain.ProductCheckout) error
	BulkCreateProductCheckout(ctx context.Context, tx *sql.Tx, productCheckout []domain.ProductCheckout) error
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

	return nil
}

func (cr *checkoutRepository) GetCheckoutHistory(ctx context.Context, db *sql.DB, queryParams domain.CheckoutHistoryQueryParams) ([]domain.GetCheckoutHistory, error) {
	var queryCondition string
	var limitOffsetClause []string
	var whereClause []string
	var orderClause []string
	var args []any

	if len(queryParams.Limit) == 0 {
		limitOffsetClause = append(limitOffsetClause, "limit 5")
	} else {
		limitOffsetClause = append(limitOffsetClause, fmt.Sprintf("limit $%d", len(args)+1))
		args = append(args, queryParams.Limit)
	}
	if len(queryParams.Offset) == 0 {
		limitOffsetClause = append(limitOffsetClause, "offset 0")
	} else {
		limitOffsetClause = append(limitOffsetClause, fmt.Sprintf("offset $%d", len(args)+1))
		args = append(args, queryParams.Offset)
	}

	val := reflect.ValueOf(queryParams)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		key := strings.ToLower(typ.Field(i).Name)
		value := val.Field(i).String()

		if key == "limit" || key == "offset" {
			continue
		}

		if len(value) < 1 {
			if key == "createdat" {
				orderClause = append(orderClause, "created_at desc")
			}
			continue
		}

		if key == "createdat" {
			if value != "asc" && value != "desc" {
				continue
			}
			key = "created_at"

			orderClause = append(orderClause, fmt.Sprintf("%s %s", key, value))
			continue
		}
	}
	if len(whereClause) > 0 {
		queryCondition += "\nAND " + strings.Join(whereClause, " AND ")
	}
	if len(orderClause) > 0 {
		queryCondition += "\nORDER BY " + strings.Join(orderClause, ", ")
	}
	queryCondition += "\n" + strings.Join(limitOffsetClause, " ")

	query := `
		SELECT c.id, c.created_at, c.user_customer_id, pc.product_id, pc.quantity, c.paid, c.change
		FROM checkouts c
		INNER JOIN product_checkouts pc ON pc.checkout_id = c.id
	`
	query += queryCondition

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	checkouts := []domain.GetCheckoutHistory{}
	for rows.Next() {
		checkout := domain.GetCheckoutHistory{}

		err := rows.Scan(&checkout.TransactionID, &checkout.CreatedAt, &checkout.CustomerID, &checkout.ProductID, &checkout.Quantity, &checkout.Paid, &checkout.Change)
		if err != nil {
			return nil, err
		}

		checkouts = append(checkouts, checkout)
	}

	return checkouts, nil
}

func (cr *checkoutRepository) BulkCreateProductCheckout(ctx context.Context, tx *sql.Tx, productCheckouts []domain.ProductCheckout) error {
	inserts := []string{}
	args := []any{}
	query := ""

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

	_, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (cr *checkoutRepository) CreateProductCheckout(ctx context.Context, tx *sql.Tx, productCheckout domain.ProductCheckout) error {
	query := `
		INSERT INTO product_checkouts (id, product_id, quantity, checkout_id)
		VALUES ($1, $2, $3, $4)
	`
	_, err := tx.ExecContext(ctx, query, productCheckout.ID, productCheckout.ProductID, productCheckout.Quantity, productCheckout.CheckoutID)
	if err != nil {
		return err
	}

	return nil
}
