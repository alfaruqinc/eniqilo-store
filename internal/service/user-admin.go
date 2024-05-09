package service

import (
	"context"
	"database/sql"
	"eniqilo-store/internal/domain"
	"eniqilo-store/internal/repository"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserAdminService interface {
	RegisterUserAdminService(ctx context.Context, userAdmin domain.RegisterUserAdminRequest) (*domain.UserAdminResponseWithAccessToken, domain.MessageErr)
	LoginUserAdminService(ctx context.Context, userAdmin domain.LoginUserAdmin) (*domain.UserAdminResponseWithAccessToken, domain.MessageErr)
	generateToken(userAdmin domain.UserAdmin) (string, error)
	mapUserAdminResponseWithAccessToken(userAdmin *domain.UserAdmin, token string) *domain.UserAdminResponseWithAccessToken
}

type userAdminService struct {
	db                  *sql.DB
	userAdminRepository repository.UserAdminRepository
	jwtSecret           string
	bcryptSalt          int
}

func NewUserAdminService(db *sql.DB, userAdminRepository repository.UserAdminRepository, jwtSecret string, bcryptSalt int) UserAdminService {
	return &userAdminService{
		db:                  db,
		userAdminRepository: userAdminRepository,
		jwtSecret:           jwtSecret,
		bcryptSalt:          bcryptSalt,
	}
}

func (u *userAdminService) RegisterUserAdminService(ctx context.Context, userAdminPayload domain.RegisterUserAdminRequest) (*domain.UserAdminResponseWithAccessToken, domain.MessageErr) {
	tx, err := u.db.Begin()
	if err != nil {
		return nil, domain.NewInternalServerError(err.Error())
	}
	defer tx.Rollback()

	phoneNumberExists, err := u.userAdminRepository.CheckPhoneNumberExists(ctx, tx, userAdminPayload.PhoneNumber)
	if err != nil {
		return nil, domain.NewInternalServerError(err.Error())
	}
	if phoneNumberExists {
		return nil, domain.NewConflictError("phone number already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userAdminPayload.Password), u.bcryptSalt)
	if err != nil {
		return nil, domain.NewInternalServerError(err.Error())
	}

	userAdmin := userAdminPayload.NewUserAdminFromDTO()
	userAdmin.Password = string(hashedPassword)

	err = u.userAdminRepository.CreateUserAdminRepository(ctx, tx, userAdmin)
	if err != nil {
		return nil, domain.NewInternalServerError(err.Error())
	}

	err = tx.Commit()
	if err != nil {
		return nil, domain.NewInternalServerError(err.Error())
	}

	token, err := u.generateToken(userAdmin)
	if err != nil {
		return nil, domain.NewInternalServerError(err.Error())
	}

	return u.mapUserAdminResponseWithAccessToken(&userAdmin, token), nil
}

func (u *userAdminService) LoginUserAdminService(ctx context.Context, userAdminPayload domain.LoginUserAdmin) (*domain.UserAdminResponseWithAccessToken, domain.MessageErr) {
	userAdmin, err := u.userAdminRepository.GetUserByPhoneNumberRepository(ctx, u.db, userAdminPayload.PhoneNumber)
	if err != nil {
		return nil, domain.NewNotFoundError("staff is not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(userAdmin.Password), []byte(userAdminPayload.Password))
	if err != nil {
		return nil, domain.NewBadRequestError("invalid password")
	}

	token, err := u.generateToken(*userAdmin)
	if err != nil {
		return nil, domain.NewInternalServerError(err.Error())
	}

	return u.mapUserAdminResponseWithAccessToken(userAdmin, token), nil
}

func (u *userAdminService) generateToken(userAdmin domain.UserAdmin) (string, error) {
	claims := jwt.MapClaims{
		"id":          userAdmin.ID,
		"phoneNumber": userAdmin.PhoneNumber,
		"exp":         time.Now().Add(time.Hour * 8).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(u.jwtSecret))

	return tokenString, err
}

func (u *userAdminService) mapUserAdminResponseWithAccessToken(userAdmin *domain.UserAdmin, token string) *domain.UserAdminResponseWithAccessToken {
	return &domain.UserAdminResponseWithAccessToken{
		ID:          userAdmin.ID,
		Name:        userAdmin.Name,
		PhoneNumber: userAdmin.PhoneNumber,
		AccessToken: token,
	}
}
