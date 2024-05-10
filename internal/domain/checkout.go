package domain

import (
	"time"

	"github.com/google/uuid"
)

type Checkout struct {
	ID             string    `db:"id"`
	CreatedAt      time.Time `db:"created_at"`
	UserCustomerID string    `db:"user_customer_id"`
	Paid           int       `db:"paid"`
	Change         *int      `db:"change"`
}

type ProductCheckout struct {
	ID         string `db:"id"`
	ProductID  string `db:"product_id"`
	Quantity   int    `db:"quantity"`
	CheckoutID string `db:"checkout_id"`
}

type ProductCheckoutRequest struct {
	ProductID string `json:"productId" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1,number"`
}

type CheckoutRequest struct {
	CustomerID     string                   `json:"customerId" binding:"required"`
	ProductDetails []ProductCheckoutRequest `json:"productDetails" binding:"required,min=1,dive"`
	Paid           int                      `json:"paid" binding:"required,min=1,number"`
	Change         *int                     `json:"change" binding:"required,min=0,number"`
}

type CheckoutResponse struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

type GetCheckoutHistory struct {
	TransactionID string    `json:"transactionId"`
	CreatedAt     time.Time `json:"createdAt"`
	CustomerID    string    `json:"customerId"`
	ProductID     string    `json:"productId"`
	Quantity      int       `json:"quantity"`
	Paid          int       `json:"paid"`
	Change        int       `json:"change"`
}

type ProductCheckoutResponse struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

type GetCheckoutHistoryResponse struct {
	TransactionID  string                    `json:"transactionId"`
	CreatedAt      time.Time                 `json:"createdAt"`
	CustomerID     string                    `json:"customerId"`
	ProductDetails []ProductCheckoutResponse `json:"productDetails" db:"product_details"`
	Paid           int                       `json:"paid"`
	Change         int                       `json:"change"`
}

type CheckoutHistoryQueryParams struct {
	CustomerId string `form:"customerId"`
	Limit      string `form:"limit"`
	Offset     string `form:"offset"`
	CreatedAt  string `form:"createdAt"`
}

func (cr *CheckoutRequest) NewCheckouts() (Checkout, []ProductCheckout) {
	id := uuid.New()
	rawCreatedAt := time.Now().Format(time.RFC3339)
	createdAt, _ := time.Parse(time.RFC3339, rawCreatedAt)

	checkout := Checkout{
		ID:             id.String(),
		CreatedAt:      createdAt,
		UserCustomerID: cr.CustomerID,
		Paid:           cr.Paid,
		Change:         cr.Change,
	}

	var productCheckouts []ProductCheckout
	for _, v := range cr.ProductDetails {
		id := uuid.New()
		productCheckout := ProductCheckout{
			ID:         id.String(),
			ProductID:  v.ProductID,
			Quantity:   v.Quantity,
			CheckoutID: checkout.ID,
		}

		productCheckouts = append(productCheckouts, productCheckout)
	}

	return checkout, productCheckouts
}
