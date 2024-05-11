package auth

import (
	"database/sql"
	"eniqilo-store/internal/domain"
	"eniqilo-store/internal/repository"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type AuthMiddleware interface {
	Authentication() gin.HandlerFunc
	validateToken(userAdmin *domain.UserAdmin, bearerToken string) error
	bindTokenToUserEntity(userAdmin *domain.UserAdmin, claim jwt.MapClaims) domain.MessageErr
	parseToken(tokenString string) (*jwt.Token, domain.MessageErr)
}

type authMiddleware struct {
	db                  *sql.DB
	jwtSecret           string
	userAdminRepository repository.UserAdminRepository
}

func NewAuthMiddleware(db *sql.DB, jwtSecret string, userAdminRepository repository.UserAdminRepository) AuthMiddleware {
	return &authMiddleware{
		db:                  db,
		jwtSecret:           jwtSecret,
		userAdminRepository: userAdminRepository,
	}
}

func (a *authMiddleware) Authentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.FullPath()
		method := ctx.Request.Method
		if method == "GET" && path == "/v1/product/customer" {
			return
		}

		invalidTokenErr := domain.NewUnauthenticatedError("invalid token")
		bearerToken := ctx.GetHeader("Authorization")

		user := domain.UserAdmin{}

		err := a.validateToken(&user, bearerToken)
		if err != nil {
			ctx.AbortWithStatusJSON(invalidTokenErr.Status(), invalidTokenErr)
			return
		}

		_, err = a.userAdminRepository.GetUserByPhoneNumberRepository(ctx, a.db, user.PhoneNumber)
		if err != nil {
			ctx.AbortWithStatusJSON(invalidTokenErr.Status(), invalidTokenErr)
			return
		}

		ctx.Set("userData", user)
		ctx.Next()
	}
}

func (a *authMiddleware) validateToken(userAdmin *domain.UserAdmin, bearerToken string) error {
	isBearer := strings.HasPrefix(bearerToken, "Bearer")
	if !isBearer {
		return domain.NewUnauthenticatedError("token should be Bearer")
	}

	splitToken := strings.Fields(bearerToken)
	if len(splitToken) != 2 {
		return domain.NewUnauthenticatedError("invalid token")
	}

	tokenString := splitToken[1]

	token, err := a.parseToken(tokenString)
	if err != nil {
		return err
	}

	var mapClaims jwt.MapClaims

	if claims, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		return domain.NewUnauthenticatedError("invalid token")
	} else {
		mapClaims = claims
	}
	err = a.bindTokenToUserEntity(userAdmin, mapClaims)
	if err != nil {
		return err
	}

	return nil
}

func (a *authMiddleware) parseToken(tokenString string) (*jwt.Token, domain.MessageErr) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.NewUnauthenticatedError("invalid token")
		}

		return []byte(a.jwtSecret), nil
	})
	if err != nil {
		return nil, domain.NewUnauthenticatedError("invalid token")
	}

	return token, nil
}

func (a *authMiddleware) bindTokenToUserEntity(userAdmin *domain.UserAdmin, claim jwt.MapClaims) domain.MessageErr {
	idString, ok := claim["id"].(string)
	if !ok {
		return domain.NewUnauthenticatedError("invalid token")
	}

	id, err := uuid.Parse(idString)
	if err != nil {
		return domain.NewUnauthenticatedError("invalid token")
	}
	userAdmin.ID = id.String()

	phoneNumber, ok := claim["phoneNumber"].(string)
	if !ok {
		return domain.NewUnauthenticatedError("invalid token")
	}
	userAdmin.PhoneNumber = phoneNumber

	return nil
}
