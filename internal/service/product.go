package service

import (
	"context"
	"database/sql"
	"eniqilo-store/internal/domain"
	"eniqilo-store/internal/repository"
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type ProductService interface {
	CreateProduct(ctx context.Context, product domain.Product) domain.MessageErr
	GetProducts(ctx context.Context, queryParams domain.ProductQueryParams) ([]domain.ProductResponse, domain.MessageErr)
	GetProductsForCustomer(ctx context.Context, queryParams domain.ProductForCustomerQueryParams) ([]domain.ProductForCustomerResponse, domain.MessageErr)
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
	err := ps.productRepository.CreateProduct(ctx, ps.db, product)
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

	limit := "5"
	qlimit, _ := strconv.Atoi(queryParams.Limit)
	if qlimit > 0 {
		limit = queryParams.Limit
	}

	offset := "0"
	qoffset, _ := strconv.Atoi(queryParams.Offset)
	if qoffset > 0 {
		qoffset = (qoffset - 1) * qlimit
		offset = strconv.Itoa(qoffset)
	}

	limitOffsetClause = append(limitOffsetClause, "limit $1 offset $2")
	args = append(args, limit, offset)

	val := reflect.ValueOf(queryParams)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		key := strings.ToLower(typ.Field(i).Name)
		value := val.Field(i).String()
		argPos := len(args) + 1

		if key == "limit" || key == "offset" {
			continue
		}

		if len(value) < 1 {
			// default order by created_at desc
			if key == "createdat" {
				orderClause = append(orderClause, "created_at desc")
				continue
			}

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
		query += "\nORDER BY " + strings.Join(orderClause, ", ") + ", sid desc"
	}
	query += "\n" + strings.Join(limitOffsetClause, " ")

	products, err := ps.productRepository.GetProducts(ctx, ps.db, query, args)
	if err != nil {
		return nil, domain.NewInternalServerError(err.Error())
	}

	return products, nil
}

func (ps *productService) GetProductsForCustomer(ctx context.Context, queryParams domain.ProductForCustomerQueryParams) ([]domain.ProductForCustomerResponse, domain.MessageErr) {
	products, err := ps.productRepository.GetProductsForCustomer(ctx, ps.db, queryParams)
	if err != nil {
		return nil, domain.NewInternalServerError(err.Error())
	}

	return products, nil
}

func (ps *productService) UpdateProductByID(ctx context.Context, product domain.Product) domain.MessageErr {
	affRow, err := ps.productRepository.UpdateProductByID(ctx, ps.db, product)
	if err != nil {
		return domain.NewInternalServerError(err.Error())
	}
	if affRow == 0 {
		return domain.NewNotFoundError("product is not found")
	}

	return nil
}

func (ps *productService) DeleteProductByID(ctx context.Context, productId string) domain.MessageErr {
	affRow, err := ps.productRepository.DeleteProductByID(ctx, ps.db, productId)
	if err != nil {
		return domain.NewInternalServerError(err.Error())
	}
	if affRow == 0 {
		return domain.NewNotFoundError("product is not found")
	}

	return nil
}
