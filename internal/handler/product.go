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
		product := domain.ProductRequest{}
		if err := ctx.ShouldBindJSON(&product); err != nil {
			err := helper.ValidateRequest(err)
			ctx.JSON(err.Status(), err)
			return
		}

		ctx.JSON(http.StatusCreated, domain.NewMessageSuccess("success create product", ""))
	}
}
