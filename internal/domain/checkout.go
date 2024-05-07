package domain

import (
	"time"

	"github.com/google/uuid"
)

type Checkout struct {
	ID             string    `db:"id"`
	CreatedAt      time.Time `db:"created_at"`
	UserCustomerID string    `db:"user_customer_id"`
	ProductID      string    `db:"product_id"`
	Paid           int       `db:"paid"`
	Change         int       `db:"change"`
}

type ProductCheckoutRequest struct {
	ProductID string `json:"productId"`
	Amount    int    `json:"amount"`
}

type CheckoutRequest struct {
	CustomerID     string                   `json:"customerId"`
	ProductDetails []ProductCheckoutRequest `json:"productDetails"`
	Paid           int                      `json:"paid"`
	Change         int                      `json:"change"`
}

type CheckoutResponse struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

func (cr *CheckoutRequest) NewCheckouts() []Checkout {
	id := uuid.New()
	rawCreatedAt := time.Now().Format(time.RFC3339)
	createdAt, _ := time.Parse(time.RFC3339, rawCreatedAt)

	var checkouts []Checkout
	for _, v := range cr.ProductDetails {
		checkout := Checkout{
			ID:             id.String(),
			CreatedAt:      createdAt,
			UserCustomerID: cr.CustomerID,
			ProductID:      v.ProductID,
			Paid:           cr.Paid,
			Change:         cr.Change,
		}

		checkouts = append(checkouts, checkout)
	}

	return checkouts
}
