package domain

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

type RegisterUserAdmin struct {
	Name        string `json:"name" binding:"required,gte=5,lte=50"`
	PhoneNumber string `json:"phoneNumber" binding:"required,gte=10,lte=16,e164"`
	Password    string `json:"password" binding:"required,gte=5,lte=15"`
}

type LoginUserAdmin struct {
	PhoneNumber string `json:"phoneNumber" binding:"required,gte=10,lte=16,e164"`
	Password    string `json:"password" binding:"required,gte=5,lte=15"`
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
