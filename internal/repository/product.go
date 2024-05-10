package repository

import (
	"context"
	"database/sql"
	"eniqilo-store/internal/domain"
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, db *sql.DB, product domain.Product) error
	GetProducts(ctx context.Context, db *sql.DB, queryParams string, args []any) ([]domain.ProductResponse, error)
	GetProductsForCustomer(ctx context.Context, db *sql.DB, queryParams domain.ProductForCustomerQueryParams) ([]domain.ProductForCustomerResponse, error)
	GetProductStockByIDs(ctx context.Context, db *sql.DB, productIds []string) ([]domain.ProductResponse, error)
	GetProductPriceByIDs(ctx context.Context, db *sql.DB, productIds []string) ([]domain.ProductResponse, error)
	UpdateProductByID(ctx context.Context, db *sql.DB, product domain.Product) (int64, error)
	DeleteProductByID(ctx context.Context, db *sql.DB, productId string) (int64, error)
	CheckProductExistsByID(ctx context.Context, db *sql.DB, productId string) (bool, error)
	CheckProductExists(ctx context.Context, db *sql.DB, IDs []string) (bool, error)
	CheckProductAvailabilities(ctx context.Context, db *sql.DB, productIDs []string) (bool, error)
	CheckProductPrice(ctx context.Context, db *sql.DB, productCheckouts []domain.ProductCheckoutRequest) (int, error)
	UpdateProductStockByID(ctx context.Context, tx *sql.Tx, product string, quantity int) error
}

type productRepository struct{}

func NewProductRepository() ProductRepository {
	return &productRepository{}
}

func (pr *productRepository) CreateProduct(ctx context.Context, db *sql.DB, product domain.Product) error {
	query := `
		INSERT INTO products (id, created_at, name, sku, category, image_url, notes, price, stock, location, is_available)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := db.ExecContext(ctx, query,
		product.ID, product.CreatedAt, product.Name, product.Sku, product.Category,
		product.ImageUrl, product.Notes, product.Price, product.Stock, product.Location,
		product.IsAvailable,
	)
	if err != nil {
		return err
	}

	return nil
}

func (pr *productRepository) GetProducts(ctx context.Context, db *sql.DB, queryParams string, args []any) ([]domain.ProductResponse, error) {
	query := `
		SELECT id, created_at, name, sku, category, image_url, stock, notes, 
				price, location, is_available
		FROM products
	`
	query += queryParams

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []domain.ProductResponse{}
	for rows.Next() {
		product := domain.ProductResponse{}

		err := rows.Scan(
			&product.ID, &product.CreatedAt, &product.Name, &product.Sku,
			&product.Category, &product.ImageUrl, &product.Stock, &product.Notes,
			&product.Price, &product.Location, &product.IsAvailable,
		)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

func (pr *productRepository) GetProductsForCustomer(ctx context.Context, db *sql.DB, queryParams domain.ProductForCustomerQueryParams) ([]domain.ProductForCustomerResponse, error) {
	var queryCondition string
	var limitOffsetClause []string
	var whereClause []string
	var orderClause []string = []string{"created_at desc"}
	var args []any

	val := reflect.ValueOf(queryParams)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		key := strings.ToLower(typ.Field(i).Name)
		value := val.Field(i).String()
		argPos := len(args) + 1

		if key == "limit" || key == "offset" {
			if key == "limit" && len(value) < 1 {
				value = "5"
			}
			if key == "offset" && len(value) < 1 {
				value = "0"
			}

			limitOffsetClause = append(limitOffsetClause, fmt.Sprintf("%s $%d", key, argPos))
			args = append(args, value)
			continue
		}

		if len(value) < 1 {
			continue
		}

		if key == "name" {
			whereClause = append(whereClause, fmt.Sprintf("%s ILIKE $%d", key, argPos))
			args = append(args, "%"+value+"%")
			continue
		}

		if key == "category" {
			if !slices.Contains(domain.ProductCategory, value) {
				continue
			}
		}

		if key == "price" {
			if value != "asc" && value != "desc" {
				continue
			}

			orderClause = append(orderClause, fmt.Sprintf("%s %s", key, value))
			continue
		}

		if key == "instock" {
			key = "stock"
			if value == "true" {
				whereClause = append(whereClause, fmt.Sprintf("%s > 0", key))
			} else if value == "false" {
				whereClause = append(whereClause, fmt.Sprintf("%s < 1", key))
			}

			continue
		}

		whereClause = append(whereClause, fmt.Sprintf("%s = $%d", key, argPos))
		args = append(args, value)
	}

	if len(whereClause) > 0 {
		queryCondition += "\nAND " + strings.Join(whereClause, " AND ")
	}
	if len(orderClause) > 0 {
		queryCondition += "\nORDER BY " + strings.Join(orderClause, ", ")
	}
	queryCondition += "\n" + strings.Join(limitOffsetClause, " ")

	query := `
		SELECT id, created_at, name, sku, category, image_url, stock, notes, 
				price, location
		FROM products
		WHERE is_available = true
	`
	query += queryCondition
	fmt.Println(query)

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []domain.ProductForCustomerResponse{}
	for rows.Next() {
		product := domain.ProductForCustomerResponse{}

		err := rows.Scan(
			&product.ID, &product.CreatedAt, &product.Name, &product.Sku,
			&product.Category, &product.ImageUrl, &product.Stock, &product.Notes,
			&product.Price, &product.Location,
		)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

func (pr *productRepository) GetProductStockByIDs(ctx context.Context, db *sql.DB, productIds []string) ([]domain.ProductResponse, error) {
	query := `
		SELECT id, name, stock
		FROM products
		WHERE id = any ($1)
	`
	rows, err := db.QueryContext(ctx, query, productIds)
	if err != nil {
		return nil, err
	}

	productsStock := []domain.ProductResponse{}
	for rows.Next() {
		productStock := domain.ProductResponse{}

		err := rows.Scan(&productStock.ID, &productStock.Name, &productStock.Stock)
		if err != nil {
			return nil, err
		}

		productsStock = append(productsStock, productStock)
	}

	return productsStock, nil
}

func (pr *productRepository) GetProductPriceByIDs(ctx context.Context, db *sql.DB, productIds []string) ([]domain.ProductResponse, error) {
	query := `
		SELECT id, price
		FROM products
		WHERE id = any ($1)
	`
	rows, err := db.QueryContext(ctx, query, productIds)
	if err != nil {
		return nil, err
	}

	productPrices := []domain.ProductResponse{}
	for rows.Next() {
		productPrice := domain.ProductResponse{}

		err := rows.Scan(&productPrice.ID, &productPrice.Price)
		if err != nil {
			return nil, err
		}

		productPrices = append(productPrices, productPrice)
	}

	return productPrices, nil
}

func (pr *productRepository) UpdateProductByID(ctx context.Context, db *sql.DB, product domain.Product) (int64, error) {
	query := `
		UPDATE products
		SET name = $2,
			sku = $3,
			category = $4,
			notes = $5,
			image_url = $6,
			price = $7,
			stock = $8,
			location = $9,
			is_available = $10
		WHERE id = $1
	`
	res, err := db.ExecContext(ctx, query,
		product.ID, product.Name, product.Sku, product.Category, product.ImageUrl,
		product.Notes, product.Price, product.Stock, product.Location, product.IsAvailable,
	)
	if err != nil {
		return 0, err
	}

	affRow, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affRow, nil
}

func (pr *productRepository) CheckProductExistsByID(ctx context.Context, db *sql.DB, productId string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM products
			WHERE id = $1
		)
	`
	var exists bool
	err := db.QueryRowContext(ctx, query, productId).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (pr *productRepository) DeleteProductByID(ctx context.Context, db *sql.DB, productId string) (int64, error) {
	query := `
		DELETE FROM products
		WHERE id = $1
	`
	res, err := db.ExecContext(ctx, query, productId)
	if err != nil {
		if err, ok := err.(*pgconn.PgError); ok {
			if err.Code == "22P02" {
				return 0, nil
			}
		}
		return 0, err
	}

	affRow, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affRow, nil
}

func (pr *productRepository) CheckProductExists(ctx context.Context, db *sql.DB, IDs []string) (bool, error) {
	query := `
		SELECT COUNT(id) = $1
		FROM products
		WHERE id = any ($2)
	`
	var exists bool
	err := db.QueryRowContext(ctx, query, len(IDs), IDs).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (cr *productRepository) CheckProductPrice(ctx context.Context, db *sql.DB, productCheckouts []domain.ProductCheckoutRequest) (int, error) {
	args := []any{}
	query := `
		SELECT 
	`

	var totalPrice int
	err := db.QueryRowContext(ctx, query, args...).Scan(&totalPrice)
	if err != nil {
		return 0, err
	}

	return totalPrice, nil
}

func (pr *productRepository) CheckProductAvailabilities(ctx context.Context, db *sql.DB, productIDs []string) (bool, error) {
	query := `
		SELECT COUNT(id) = $1
		FROM products
		WHERE id = any ($2)
			AND is_available = true
	`
	var exists bool
	err := db.QueryRowContext(ctx, query, len(productIDs), productIDs).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (pr *productRepository) UpdateProductStockByID(ctx context.Context, tx *sql.Tx, id string, quantity int) error {
	query := `
		UPDATE products
		SET stock = stock - $1
		WHERE id = $2
	`
	_, err := tx.ExecContext(ctx, query, quantity, id)
	if err != nil {
		return err
	}

	return nil
}
