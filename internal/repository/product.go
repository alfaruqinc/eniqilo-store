package repository

import (
	"context"
	"database/sql"
	"eniqilo-store/internal/domain"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, tx *sql.Tx, product domain.Product) error
	UpdateProduct(ctx context.Context, tx *sql.Tx, product domain.Product) error
}

type productRepository struct{}

func NewProductRepository() ProductRepository {
	return &productRepository{}
}

func (pr *productRepository) CreateProduct(ctx context.Context, tx *sql.Tx, product domain.Product) error {
	query := `
		INSERT INTO products (id, created_at, name, sku, category, image_url, notes, price, stock, location, is_available)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := tx.ExecContext(ctx, query,
		product.ID, product.CreatedAt, product.Name, product.Sku, product.Category,
		product.ImageUrl, product.Notes, product.Price, product.Stock, product.Location,
		product.IsAvailable,
	)
	if err != nil {
		return err
	}

	return nil
}

func (pr *productRepository) UpdateProduct(ctx context.Context, tx *sql.Tx, product domain.Product) error {
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
	_, err := tx.ExecContext(ctx, query,
		product.ID, product.Name, product.Sku, product.Category, product.ImageUrl,
		product.Notes, product.Price, product.Stock, product.Location, product.IsAvailable,
	)
	if err != nil {
		return err
	}

	return nil
}
