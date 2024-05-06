package server

import (
	"net/http"
	"project-sprint-w2/internal/handler"
	"project-sprint-w2/internal/repository"
	"project-sprint-w2/internal/service"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	db := s.db.GetDB()

	userAdminRepository := repository.NewUserAdminRepository()

	userAdminService := service.NewUserAdminService(db, userAdminRepository)

	userAdminHandler := handler.NewUserAdminHandler(userAdminService)

	r := gin.Default()

	r.GET("/", s.HelloWorldHandler)

	r.GET("/health", s.healthHandler)

	staff := r.Group("/staff")
	staff.POST("/register", userAdminHandler.RegisterUserAdminHandler())
	staff.POST("/login", userAdminHandler.LoginUserAdminHandler())

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
