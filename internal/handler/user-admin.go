package handler

import (
	"eniqilo-store/internal/domain"
	"eniqilo-store/internal/helper"
	"eniqilo-store/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
		userAdmin := domain.RegisterUserAdmin{}
		if err := c.ShouldBindJSON(&userAdmin); err != nil {
			if err, ok := err.(validator.ValidationErrors); ok {
				for _, fe := range err {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": helper.MsgForTag(fe),
					})
					return
				}
			}

			c.JSON(http.StatusInternalServerError, err)
		}

		response, err := u.userAdminService.RegisterUserAdminService(c.Request.Context(), userAdmin)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "Register User Admin",
			"data":    response,
		})
	}
}

func (u *userAdminHandler) LoginUserAdminHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Login User Admin",
		})
	}
}
