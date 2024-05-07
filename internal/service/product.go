package service

import (
	"context"
	"database/sql"
	"eniqilo-store/internal/domain"
	"eniqilo-store/internal/repository"
)

type ProductService interface {
	CreateProduct(ctx context.Context, product domain.Product) domain.MessageErr
	UpdateProductByID(ctx context.Context, product domain.Product) domain.MessageErr
	DeleteProductByID(ctx context.Context, productId string) domain.MessageErr
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
	if err != nil {
		return domain.NewInternalServerError(err.Error())
	}

	err = tx.Commit()
	if err != nil {
		return domain.NewInternalServerError(err.Error())
	}

	return nil
}

func (ps *productService) UpdateProductByID(ctx context.Context, product domain.Product) domain.MessageErr {
	tx, err := ps.db.Begin()
	if err != nil {
		return domain.NewInternalServerError(err.Error())
	}
	defer tx.Rollback()

	productExists, err := ps.productRepository.CheckProductExistsByID(ctx, tx, product.ID)
	if err != nil {
		return domain.NewInternalServerError(err.Error())
	}
	if !productExists {
		return domain.NewNotFoundError("product is not found")
	}

	err = ps.productRepository.UpdateProductByID(ctx, tx, product)
	if err != nil {
		return domain.NewInternalServerError(err.Error())
	}

	err = tx.Commit()
	if err != nil {
		return domain.NewInternalServerError(err.Error())
	}

	return nil
}

func (ps *productService) DeleteProductByID(ctx context.Context, productId string) domain.MessageErr {
	tx, err := ps.db.Begin()
	if err != nil {
		return domain.NewInternalServerError(err.Error())
	}
	defer tx.Rollback()

	productExists, err := ps.productRepository.CheckProductExistsByID(ctx, tx, productId)
	if err != nil {
		return domain.NewInternalServerError(err.Error())
	}
	if !productExists {
		return domain.NewNotFoundError("product is not found")
	}

	err = ps.productRepository.DeleteProductByID(ctx, tx, productId)
	if err != nil {
		return domain.NewInternalServerError(err.Error())
	}

	err = tx.Commit()
	if err != nil {
		return domain.NewInternalServerError(err.Error())
	}

	return nil
}
