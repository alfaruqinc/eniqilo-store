package service

import (
	"context"
	"database/sql"
	"eniqilo-store/internal/domain"
	"eniqilo-store/internal/repository"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserAdminService interface {
	RegisterUserAdminService(ctx context.Context, userAdmin domain.RegisterUserAdminRequest) (*domain.UserAdminResponseWithAccessToken, error)
	LoginUserAdminService(ctx context.Context, userAdmin domain.LoginUserAdmin) (*domain.UserAdminResponseWithAccessToken, error)
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

func (u *userAdminService) RegisterUserAdminService(ctx context.Context, userAdminPayload domain.RegisterUserAdminRequest) (*domain.UserAdminResponseWithAccessToken, error) {
	tx, err := u.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userAdminPayload.Password), u.bcryptSalt)
	if err != nil {
		return nil, err
	}

	userAdmin := userAdminPayload.NewUserAdminFromDTO()
	userAdmin.Password = string(hashedPassword)

	err = u.userAdminRepository.CreateUserAdminRepository(ctx, tx, userAdmin)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	token, err := u.generateToken(userAdmin)
	if err != nil {
		return nil, err
	}

	return u.mapUserAdminResponseWithAccessToken(&userAdmin, token), nil
}

func (u *userAdminService) LoginUserAdminService(ctx context.Context, userAdminPayload domain.LoginUserAdmin) (*domain.UserAdminResponseWithAccessToken, error) {
	tx, err := u.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	userAdmin, err := u.userAdminRepository.GetUserByPhoneNumberRepository(ctx, tx, userAdminPayload.PhoneNumber)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(userAdmin.Password), []byte(userAdminPayload.Password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	token, err := u.generateToken(*userAdmin)
	if err != nil {
		return nil, err
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
	tokenString, _ := token.SignedString([]byte(u.jwtSecret))

	return tokenString, nil
}

func (u *userAdminService) mapUserAdminResponseWithAccessToken(userAdmin *domain.UserAdmin, token string) *domain.UserAdminResponseWithAccessToken {
	return &domain.UserAdminResponseWithAccessToken{
		ID:          userAdmin.ID,
		Name:        userAdmin.Name,
		PhoneNumber: userAdmin.PhoneNumber,
		AccessToken: token,
	}
}
