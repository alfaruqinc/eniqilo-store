package service

import (
	"context"
	"database/sql"
	"eniqilo-store/internal/domain"
	"eniqilo-store/internal/repository"
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/google/uuid"
)

type ProductService interface {
	CreateProduct(ctx context.Context, product domain.Product) domain.MessageErr
	GetProducts(ctx context.Context, queryParams domain.ProductQueryParams) ([]domain.ProductResponse, domain.MessageErr)
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

func (ps *productService) GetProducts(ctx context.Context, queryParams domain.ProductQueryParams) ([]domain.ProductResponse, domain.MessageErr) {
	var query string
	var limitOffsetClause []string
	var whereClause []string
	var orderClause []string
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

		if key == "id" {
			if _, err := uuid.Parse(value); err != nil {
				continue
			}
		}

		if key == "name" {
			whereClause = append(whereClause, fmt.Sprintf("%s ILIKE $%d", key, argPos))
			args = append(args, "%"+value+"%")
			continue
		}

		if key == "isavailable" {
			key = "is_available"
		}

		if key == "category" {
			if !slices.Contains(domain.ProductCategory, value) {
				continue
			}
		}

		if key == "price" || key == "createdat" {
			if value != "asc" && value != "desc" {
				continue
			}
			if key == "createdat" {
				key = "created_at"
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
		query += "\nWHERE " + strings.Join(whereClause, " AND ")
	}
	if len(orderClause) > 0 {
		query += "\nORDER BY " + strings.Join(orderClause, ", ")
	}
	query += "\n" + strings.Join(limitOffsetClause, " ")

	products, err := ps.productRepository.GetProducts(ctx, ps.db, query, args)
	if err != nil {
		return nil, domain.NewInternalServerError(err.Error())
	}

	return products, nil
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
