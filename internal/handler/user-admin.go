package handler

import (
	"project-sprint-w2/internal/service"

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
		c.JSON(200, gin.H{
			"message": "Register User Admin",
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
