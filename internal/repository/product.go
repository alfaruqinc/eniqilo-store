package repository

import (
	"context"
	"database/sql"
	"eniqilo-store/internal/domain"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, db *sql.DB, product domain.Product) error
	GetProducts(ctx context.Context, db *sql.DB, queryParams string, args []any) ([]domain.ProductResponse, error)
	GetProductsForCustomer(ctx context.Context, db *sql.DB, queryParams string, args []any) ([]domain.ProductForCustomerResponse, error)
	UpdateProductByID(ctx context.Context, db *sql.DB, product domain.Product) (int64, error)
	DeleteProductByID(ctx context.Context, db *sql.DB, productId string) (int64, error)
	CheckProductExistsByID(ctx context.Context, db *sql.DB, productId string) (bool, error)
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

func (pr *productRepository) GetProductsForCustomer(ctx context.Context, db *sql.DB, queryParams string, args []any) ([]domain.ProductForCustomerResponse, error) {
	query := `
		SELECT id, created_at, name, sku, category, image_url, stock, notes, 
				price, location
		FROM products
	`
	query += queryParams

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
		return 0, err
	}

	affRow, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affRow, nil
}
