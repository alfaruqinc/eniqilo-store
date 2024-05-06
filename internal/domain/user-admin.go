package domain

////
// Vars
////

var (
	UserAdminRoleStaff = "staff"
)

////
// Structs
////

////
// DTO
////

type RegisterUserAdmin struct {
	Name        string `json:"name" validate:"required,phoneNumber"`
	PhoneNumber string `json:"phone_number" validate:"required,min=5,max=5s"`
	Password    string `json:"password" validate:"required,min=5,max=15"`
}

type LoginUserAdmin struct {
	PhoneNumber string `json:"phone_number" validate:"required,min=5,max=5s"`
	Password    string `json:"password" validate:"required,min=5,max=15"`
}

type UserAdmin struct {
	ID          string `db:"id"`
	Name        string `db:"name"`
	PhoneNumber string `db:"phone_number"`
	Role        string `db:"role"`
	Password    string `db:"password"`
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
