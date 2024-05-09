package handler

import (
	"eniqilo-store/internal/domain"
	"eniqilo-store/internal/helper"
	"eniqilo-store/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ProductHandler interface {
	CreateProduct() gin.HandlerFunc
	GetProducts() gin.HandlerFunc
	GetProductsForCustomer() gin.HandlerFunc
	UpdateProductByID() gin.HandlerFunc
	DeleteProductByID() gin.HandlerFunc
}

type productHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) ProductHandler {
	return &productHandler{
		productService: productService,
	}
}

func (ph *productHandler) CreateProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		productBody := domain.ProductRequest{}
		if err := ctx.ShouldBindJSON(&productBody); err != nil {
			err := helper.ValidateRequest(err)
			ctx.JSON(err.Status(), err)
			return
		}

		product := productBody.NewProduct()
		err := ph.productService.CreateProduct(ctx, product)
		if err != nil {
			ctx.JSON(err.Status(), err)
			return
		}

		productResponse := domain.CreateProductResponse{
			ID:        product.ID,
			CreatedAt: product.CreatedAt,
		}

		ctx.JSON(http.StatusCreated, domain.NewMessageSuccess("success create product", productResponse))
	}
}

func (ph *productHandler) GetProducts() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var queryParams domain.ProductQueryParams
		ctx.ShouldBindQuery(&queryParams)

		products, err := ph.productService.GetProducts(ctx, queryParams)
		if err != nil {
			err, _ := err.(domain.MessageErr)
			ctx.JSON(err.Status(), err)
			return
		}

		ctx.JSON(http.StatusOK, domain.NewMessageSuccess("success get products", products))
	}
}

func (ph *productHandler) GetProductsForCustomer() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var queryParams domain.ProductForCustomerQueryParams
		ctx.ShouldBindQuery(&queryParams)

		products, err := ph.productService.GetProductsForCustomer(ctx, queryParams)
		if err != nil {
			err, _ := err.(domain.MessageErr)
			ctx.JSON(err.Status(), err)
			return
		}

		ctx.JSON(http.StatusOK, domain.NewMessageSuccess("success get products for customer", products))
	}
}

func (ph *productHandler) UpdateProductByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		productBody := domain.ProductRequest{}
		if err := ctx.ShouldBindJSON(&productBody); err != nil {
			err := helper.ValidateRequest(err)
			ctx.JSON(err.Status(), err)
			return
		}

		productId := ctx.Param("id")

		product := productBody.NewProduct()
		product.ID = productId
		err := ph.productService.UpdateProductByID(ctx, product)
		if err != nil {
			ctx.JSON(err.Status(), err)
			return
		}

		productResponse := domain.UpdateProductResponse{
			ID:        product.ID,
			UpdatedAt: product.CreatedAt,
		}

		ctx.JSON(http.StatusCreated, domain.NewMessageSuccess("success update product", productResponse))
	}
}

func (ph *productHandler) DeleteProductByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		productId := ctx.Param("id")

		err := ph.productService.DeleteProductByID(ctx, productId)
		if err != nil {
			ctx.JSON(err.Status(), err)
			return
		}

		rawDeletedAt := time.Now().Format(time.RFC3339)
		deletedAt, _ := time.Parse(time.RFC3339, rawDeletedAt)
		productResponse := domain.DeleteProductResponse{
			ID:        productId,
			DeletedAt: deletedAt,
		}

		ctx.JSON(http.StatusCreated, domain.NewMessageSuccess("success delete product", productResponse))
	}
}
