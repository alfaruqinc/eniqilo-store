package service

import (
	"context"
	"database/sql"
	"eniqilo-store/internal/domain"
	"eniqilo-store/internal/repository"
)

type CheckoutService interface {
	CreateCheckout(ctx context.Context, body domain.CheckoutRequest) domain.MessageErr
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
		return domain.NewInternalServerError(err.Error())
	}

	err = tx.Commit()
	if err != nil {
		return domain.NewInternalServerError(err.Error())
	}

	return nil
}
