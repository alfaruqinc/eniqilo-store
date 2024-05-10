package server

import (
	"eniqilo-store/internal/auth"
	"eniqilo-store/internal/handler"
	"eniqilo-store/internal/repository"
	"eniqilo-store/internal/service"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var (
	jwtSecret     = os.Getenv("JWT_SECRET")
	bcryptSalt, _ = strconv.Atoi(os.Getenv("BCRYPT_SALT"))
)

func (s *Server) RegisterRoutes() http.Handler {
	db := s.db.GetDB()

	userAdminRepository := repository.NewUserAdminRepository()
	productRepository := repository.NewProductRepository()
	userCustomerRepository := repository.NewUserCustomerRepository()
	checkoutRepository := repository.NewCheckoutRepository()

	userAdminService := service.NewUserAdminService(db, userAdminRepository, jwtSecret, bcryptSalt)
	productService := service.NewProductService(db, productRepository)
	userCustomerService := service.NewUserCustomerService(db, userCustomerRepository)
	checkoutService := service.NewCheckoutService(db, checkoutRepository, userCustomerRepository, productRepository)
	auths := auth.NewAuthMiddleware(db, jwtSecret, userAdminRepository)

	userAdminHandler := handler.NewUserAdminHandler(userAdminService)
	productHandler := handler.NewProductHandler(productService)
	userCustomerHandler := handler.NewUserCustomerHandler(userCustomerService)
	checkoutHandler := handler.NewCheckoutHandler(checkoutService)

	r := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("validurl", validURL)
		v.RegisterValidation("phonenumber", validPhonenumber)
	}

	r.GET("/", s.HelloWorldHandler)

	r.GET("/health", s.healthHandler)

	apiV1 := r.Group("/v1")

	staff := apiV1.Group("/staff")
	staff.POST("/register", userAdminHandler.RegisterUserAdminHandler())
	staff.POST("/login", userAdminHandler.LoginUserAdminHandler())

	product := apiV1.Group("/product")
	product.Use(auths.Authentication())
	product.POST("", productHandler.CreateProduct())
	product.GET("", productHandler.GetProducts())
	product.PUT(":id", productHandler.UpdateProductByID())
	product.DELETE(":id", productHandler.DeleteProductByID())
	product.GET("/customer", productHandler.GetProductsForCustomer())

	checkout := product.Group("/checkout")
	checkout.POST("", checkoutHandler.CreateCheckout())
	checkout.GET("/history", checkoutHandler.GetCheckoutHistory())

	customer := apiV1.Group("/customer")
	customer.Use(auths.Authentication())
	customer.GET("", userCustomerHandler.GetUserCustomers())
	customer.POST("/register", userCustomerHandler.CreateUserCustomer())

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

func validURL(fl validator.FieldLevel) bool {
	urlString := fl.Field().String()
	if !strings.Contains(urlString, ".") {
		return false
	}
	_, err := url.ParseRequestURI(urlString)
	if err != nil {
		return false
	}
	u, err := url.Parse(urlString)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}
	return true
}

func validPhonenumber(fl validator.FieldLevel) bool {
	phoneNumber := fl.Field().String()
	if !strings.HasPrefix(phoneNumber, "+") {
		return false
	}

	return true
}
