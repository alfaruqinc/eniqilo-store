package handler

import (
	"eniqilo-store/internal/domain"
	"eniqilo-store/internal/helper"
	"eniqilo-store/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserCustomerHandler interface {
	CreateUserCustomer() gin.HandlerFunc
}

type userCustomerHandler struct {
	userCustomerSerivce service.UserCustomerService
}

func NewUserCustomerHandler(userCustomerService service.UserCustomerService) UserCustomerHandler {
	return &userCustomerHandler{
		userCustomerSerivce: userCustomerService,
	}
}

func (uch *userCustomerHandler) CreateUserCustomer() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body domain.RegisterUserCustomerRequest
		if err := ctx.ShouldBindJSON(&body); err != nil {
			err := helper.ValidateRequest(err)
			ctx.JSON(err.Status(), err)
			return
		}

		userCustomer := body.NewUserCustomer()
		err := uch.userCustomerSerivce.CreateUserCustomer(ctx.Request.Context(), userCustomer)
		if err != nil {
			ctx.JSON(err.Status(), err)
			return
		}

		response := domain.RegisterUserCustomerResponse{
			ID:          userCustomer.ID,
			PhoneNumber: userCustomer.PhoneNumber,
			Name:        userCustomer.Name,
		}

		ctx.JSON(http.StatusCreated, domain.NewMessageSuccess("success register customer", response))
	}
}
