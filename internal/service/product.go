package service

import (
	"context"
	"database/sql"
	"eniqilo-store/internal/domain"
	"eniqilo-store/internal/repository"
)

type ProductService interface {
	CreateProduct(ctx context.Context, product domain.Product) domain.MessageErr
}

type productService struct {
	db                *sql.DB
	productRepository repository.ProductRepository
}

func NewProductService(db *sql.DB, productRepository repository.ProductRepository) ProductService {
	return &productService{
		db:                db,
		productRepository: productRepository,
	}
}

func (ps *productService) CreateProduct(ctx context.Context, product domain.Product) domain.MessageErr {
	tx, err := ps.db.Begin()
	if err != nil {
		return domain.NewInternalServerError(err.Error())
	}
	defer tx.Rollback()

	err = ps.productRepository.CreateProduct(ctx, tx, product)

	err = tx.Commit()
	if err != nil {
		return domain.NewInternalServerError(err.Error())
	}

	return nil
}
