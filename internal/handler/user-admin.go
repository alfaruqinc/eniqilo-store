package handler

import (
	"eniqilo-store/internal/domain"
	"eniqilo-store/internal/helper"
	"eniqilo-store/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserAdminHandler interface {
	RegisterUserAdminHandler() gin.HandlerFunc
	LoginUserAdminHandler() gin.HandlerFunc
}

type userAdminHandler struct {
	userAdminService service.UserAdminService
}

func NewUserAdminHandler(userAdminService service.UserAdminService) UserAdminHandler {
	return &userAdminHandler{
		userAdminService: userAdminService,
	}
}

func (u *userAdminHandler) RegisterUserAdminHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userAdmin := domain.RegisterUserAdminRequest{}
		if err := c.ShouldBindJSON(&userAdmin); err != nil {
			err := helper.ValidateRequest(err)
			c.JSON(err.Status(), err)
			return
		}

		response, err := u.userAdminService.RegisterUserAdminService(c.Request.Context(), userAdmin)
		if err != nil {
			c.JSON(err.Status(), err)
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Register User Admin",
			"data":    response,
		})
	}
}

func (u *userAdminHandler) LoginUserAdminHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userAdmin := domain.LoginUserAdmin{}
		if err := c.ShouldBindJSON(&userAdmin); err != nil {
			err := helper.ValidateRequest(err)
			c.JSON(err.Status(), err)
			return
		}

		response, err := u.userAdminService.LoginUserAdminService(c.Request.Context(), userAdmin)
		if err != nil {
			c.JSON(err.Status(), err)
			return
		}

		c.JSON(http.StatusOK, domain.NewMessageSuccess("success login", response))
	}
}
