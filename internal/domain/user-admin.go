package domain

import (
	"time"

	"github.com/google/uuid"
)

////
// Vars
////

var UserAdminRoleStaff = "staff"

////
// Structs
////

////
// DTO
////

type RegisterUserAdminRequest struct {
	Name        string `json:"name" binding:"required,gte=5,lte=50,fullname"`
	PhoneNumber string `json:"phoneNumber" binding:"required,gte=10,lte=16,e164"`
	Password    string `json:"password" binding:"required,gte=5,lte=15"`
}

type LoginUserAdmin struct {
	PhoneNumber string `json:"phoneNumber" binding:"required,gte=10,lte=16,e164"`
	Password    string `json:"password" binding:"required,gte=5,lte=15"`
}

type UserAdmin struct {
	ID          string    `db:"id"`
	CreatedAt   time.Time `db:"created_at"`
	Name        string    `db:"name"`
	PhoneNumber string    `db:"phone_number"`
	Role        string    `db:"role"`
	Password    string    `db:"password"`
}

////
// Response
////

type UserAdminResponseWithAccessToken struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
	AccessToken string `json:"accessToken"`
}

func (rua *RegisterUserAdminRequest) NewUserAdminFromDTO() UserAdmin {
	id := uuid.New()
	rawCreatedAt := time.Now().Format(time.RFC3339)
	createdAt, _ := time.Parse(time.RFC3339, rawCreatedAt)

	return UserAdmin{
		ID:          id.String(),
		CreatedAt:   createdAt,
		Name:        rua.Name,
		PhoneNumber: rua.PhoneNumber,
		Password:    rua.Password,
		Role:        UserAdminRoleStaff,
	}
}
