package handler

import (
	"eniqilo-store/internal/domain"
	"eniqilo-store/internal/helper"
	"eniqilo-store/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CheckoutHandler interface {
	CreateCheckout() gin.HandlerFunc
}

type checkoutHandler struct {
	checkoutSerivce service.CheckoutService
}

func NewCheckoutHandler(checkoutService service.CheckoutService) CheckoutHandler {
	return &checkoutHandler{
		checkoutSerivce: checkoutService,
	}
}

func (ch *checkoutHandler) CreateCheckout() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body domain.CheckoutRequest
		if err := ctx.ShouldBindJSON(&body); err != nil {
			err := helper.ValidateRequest(err)
			ctx.JSON(err.Status(), err)
			return
		}

		err := ch.checkoutSerivce.CreateCheckout(ctx.Request.Context(), body)
		if err != nil {
			ctx.JSON(err.Status(), err)
			return
		}

		ctx.JSON(http.StatusCreated, domain.NewMessageSuccess("success checkout", ""))
	}
}
