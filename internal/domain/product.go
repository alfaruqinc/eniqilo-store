package domain

import (
	"time"

	"github.com/google/uuid"
)

var (
	ProductCategoryClothing    = "Clothing"
	ProductCategoryAccessories = "Accessories"
	ProductCategoryFootwear    = "Footwear"
	ProductCategoryBeverages   = "Beverages"
)

var ProductCategory = []string{
	ProductCategoryClothing,
	ProductCategoryAccessories,
	ProductCategoryFootwear,
	ProductCategoryBeverages,
}

type Product struct {
	ID          string    `db:"id"`
	CreatedAt   time.Time `db:"created_at"`
	Name        string    `db:"name"`
	Sku         string    `db:"sku"`
	Category    string    `db:"category"`
	ImageUrl    string    `db:"image_url"`
	Notes       string    `db:"notes"`
	Price       int       `db:"price"`
	Stock       int       `db:"stock"`
	Location    string    `db:"stock"`
	IsAvailable bool      `db:"is_available"`
}

type ProductRequest struct {
	Name        string `json:"name" binding:"required,gte=1,lte=30"`
	Sku         string `json:"sku" binding:"required,gte=1,lte=30"`
	Category    string `json:"category" binding:"required,oneof=Clothing Accessories Footwear Beverages"`
	ImageUrl    string `json:"imageUrl" binding:"required,url"`
	Notes       string `json:"notes" binding:"required,gte=1,lte=200"`
	Price       int    `json:"price" binding:"min=1"`
	Stock       int    `json:"stock" binding:"required,min=0,max=100000"`
	Location    string `json:"location" binding:"required,gte=1,lte=200"`
	IsAvailable bool   `json:"isAvailable" binding:"boolean"`
}

type CreateProductResponse struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

type ProductResponse struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	Name        string    `json:"name"`
	Sku         string    `json:"sku"`
	Category    string    `json:"category"`
	ImageUrl    string    `json:"imageUrl"`
	Notes       string    `json:"notes"`
	Price       int       `json:"price"`
	Stock       int       `json:"stock"`
	Location    string    `json:"location"`
	IsAvailable bool      `json:"isAvailable"`
}

type ProductForCustomerResponse struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Name      string    `json:"name"`
	Sku       string    `json:"sku"`
	Category  string    `json:"category"`
	ImageUrl  string    `json:"imageUrl"`
	Notes     string    `json:"notes"`
	Price     int       `json:"price"`
	Stock     int       `json:"stock"`
	Location  string    `json:"location"`
}

type UpdateProductResponse struct {
	ID        string    `json:"id"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type DeleteProductResponse struct {
	ID        string    `json:"id"`
	DeletedAt time.Time `json:"deletedAt"`
}

type ProductQueryParams struct {
	Id          string `form:"id"`
	Limit       string `form:"limit"`
	Offset      string `form:"offset"`
	Name        string `form:"name"`
	IsAvailable string `form:"isAvailable"`
	Category    string `form:"category"`
	Sku         string `form:"sku"`
	Price       string `form:"price"`
	InStock     string `form:"inStock"`
	CreatedAt   string `form:"createdAt"`
}

type ProductForCustomerQueryParams struct {
	Limit    string `form:"limit"`
	Offset   string `form:"offset"`
	Name     string `form:"name"`
	Category string `form:"category"`
	Sku      string `form:"sku"`
	Price    string `form:"price"`
	InStock  string `form:"inStock"`
}

func (pr *ProductRequest) NewProduct() Product {
	id := uuid.New()
	rawCreatedAt := time.Now().Format(time.RFC3339)
	createdAt, _ := time.Parse(time.RFC3339, rawCreatedAt)

	return Product{
		ID:          id.String(),
		CreatedAt:   createdAt,
		Name:        pr.Name,
		Sku:         pr.Sku,
		Category:    pr.Category,
		ImageUrl:    pr.ImageUrl,
		Notes:       pr.Notes,
		Price:       pr.Price,
		Stock:       pr.Stock,
		Location:    pr.Location,
		IsAvailable: pr.IsAvailable,
	}
}
