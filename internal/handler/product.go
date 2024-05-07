package handler

import (
	"eniqilo-store/internal/domain"
	"eniqilo-store/internal/helper"
	"eniqilo-store/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductHandler interface {
	CreateProduct() gin.HandlerFunc
	UpdateProduct() gin.HandlerFunc
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

func (ph *productHandler) UpdateProduct() gin.HandlerFunc {
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
		err := ph.productService.UpdateProduct(ctx, product)
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
