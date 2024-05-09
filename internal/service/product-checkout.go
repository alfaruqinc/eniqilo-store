package service

import (
	"context"
	"database/sql"
	"eniqilo-store/internal/domain"
	"eniqilo-store/internal/repository"

	"github.com/jackc/pgx/v5/pgconn"
)

type CheckoutService interface {
	CreateCheckout(ctx context.Context, body domain.CheckoutRequest) domain.MessageErr
	GetCheckoutHistory(ctx context.Context, queryParams domain.CheckoutHistoryQueryParams) ([]domain.GetCheckoutHistoryResponse, domain.MessageErr)
}

type checkoutService struct {
	db                 *sql.DB
	checkoutRepository repository.CheckoutRepository
}

func NewCheckoutService(db *sql.DB, checkoutRepository repository.CheckoutRepository) CheckoutService {
	return &checkoutService{
		db:                 db,
		checkoutRepository: checkoutRepository,
	}
}

func (cs *checkoutService) CreateCheckout(ctx context.Context, body domain.CheckoutRequest) domain.MessageErr {
	checkout, productCheckouts := body.NewCheckouts()

	tx, err := cs.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.NewInternalServerError(err.Error())
	}
	defer tx.Rollback()

	err = cs.checkoutRepository.CreateCheckout(ctx, tx, checkout, productCheckouts)
	if err != nil {
		if err, ok := err.(*pgconn.PgError); ok {
			customerNotFound := err.Code == "23503" && err.ConstraintName == "fk_user_customer_id_checkouts"
			productsNotFound := err.Code == "23503" && err.ConstraintName == "fk_product_id_product_checkouts"
			if customerNotFound {
				return domain.NewNotFoundError("customer is not found")
			} else if productsNotFound {
				return domain.NewNotFoundError("one of your products is not found")
			}
		}

		return domain.NewInternalServerError(err.Error())
	}

	err = tx.Commit()
	if err != nil {
		return domain.NewInternalServerError(err.Error())
	}

	return nil
}

func (cs *checkoutService) GetCheckoutHistory(ctx context.Context, queryParams domain.CheckoutHistoryQueryParams) ([]domain.GetCheckoutHistoryResponse, domain.MessageErr) {
	checkouts, err := cs.checkoutRepository.GetCheckoutHistory(ctx, cs.db, queryParams)
	if err != nil {
		return nil, domain.NewInternalServerError(err.Error())
	}

	productDetailsMap := map[string][]domain.ProductCheckoutResponse{}
	for _, chk := range checkouts {
		productDetailsMap[chk.TransactionID] = append(productDetailsMap[chk.TransactionID], domain.ProductCheckoutResponse{
			ProductID: chk.ProductID,
			Quantity:  chk.Quantity,
		})
	}

	checkoutHistory := []domain.GetCheckoutHistoryResponse{}
	uniqueMap := map[string]bool{}
	for _, chk := range checkouts {
		if !uniqueMap[chk.TransactionID] {
			uniqueMap[chk.TransactionID] = true
			history := domain.GetCheckoutHistoryResponse{
				TransactionID:  chk.TransactionID,
				Paid:           chk.Paid,
				Change:         chk.Change,
				ProductDetails: productDetailsMap[chk.TransactionID],
			}
			checkoutHistory = append(checkoutHistory, history)
		}
	}

	return checkoutHistory, nil
}
