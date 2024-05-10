package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserCustomer struct {
	ID          string    `db:"id"`
	CreatedAt   time.Time `db:"created_at"`
	Name        string    `db:"name"`
	PhoneNumber string    `db:"phone_number"`
}

type RegisterUserCustomerRequest struct {
	Name        string `json:"name" binding:"required,gte=5,lte=50"`
	PhoneNumber string `json:"phoneNumber" binding:"required,gte=10,lte=17,phonenumber"`
}

type RegisterUserCustomerResponse struct {
	ID          string `json:"userId"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
}

type UserCustomerResponse struct {
	ID          string `json:"userId"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
}

type UserCustomerQueryParams struct {
	Name        string `form:"name"`
	PhoneNumber string `form:"phoneNumber"`
}

func (cr *RegisterUserCustomerRequest) NewUserCustomer() UserCustomer {
	id := uuid.New()
	rawCreatedAt := time.Now().Format(time.RFC3339)
	createdAt, _ := time.Parse(time.RFC3339, rawCreatedAt)

	return UserCustomer{
		ID:          id.String(),
		CreatedAt:   createdAt,
		Name:        cr.Name,
		PhoneNumber: cr.PhoneNumber,
	}
}
