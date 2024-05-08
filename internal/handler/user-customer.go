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
	GetUserCustomers() gin.HandlerFunc
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

func (uch *userCustomerHandler) GetUserCustomers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var queryParams domain.UserCustomerQueryParams
		ctx.ShouldBindQuery(&queryParams)

		customers, err := uch.userCustomerSerivce.GetUserCustomers(ctx, queryParams)
		if err != nil {
			ctx.JSON(err.Status(), err)
			return
		}

		ctx.JSON(http.StatusOK, domain.NewMessageSuccess("success get customers", customers))
	}
}
