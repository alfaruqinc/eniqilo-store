package server

import (
	"eniqilo-store/internal/handler"
	"eniqilo-store/internal/repository"
	"eniqilo-store/internal/service"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	jwtSecret     = os.Getenv("JWT_SECRET")
	bcryptSalt, _ = strconv.Atoi(os.Getenv("BCRYPT_SALT"))
)

func (s *Server) RegisterRoutes() http.Handler {
	db := s.db.GetDB()

	userAdminRepository := repository.NewUserAdminRepository()

	userAdminService := service.NewUserAdminService(db, userAdminRepository, jwtSecret, bcryptSalt)

	userAdminHandler := handler.NewUserAdminHandler(userAdminService)

	r := gin.Default()

	r.GET("/", s.HelloWorldHandler)

	r.GET("/health", s.healthHandler)

	apiV1 := r.Group("/v1")

	staff := apiV1.Group("/staff")
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
