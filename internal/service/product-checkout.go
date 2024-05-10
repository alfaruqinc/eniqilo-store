package service

import (
	"context"
	"database/sql"
	"eniqilo-store/internal/domain"
	"eniqilo-store/internal/repository"
	"fmt"
)

type CheckoutService interface {
	CreateCheckout(ctx context.Context, body domain.CheckoutRequest) domain.MessageErr
	GetCheckoutHistory(ctx context.Context, queryParams domain.CheckoutHistoryQueryParams) ([]domain.GetCheckoutHistoryResponse, domain.MessageErr)
}

type checkoutService struct {
	db                     *sql.DB
	checkoutRepository     repository.CheckoutRepository
	userCustomerRepository repository.UserCustomerRepository
	productRepository      repository.ProductRepository
}

func NewCheckoutService(db *sql.DB, checkoutRepository repository.CheckoutRepository, userCustomerRepository repository.UserCustomerRepository, productRepository repository.ProductRepository) CheckoutService {
	return &checkoutService{
		db:                     db,
		checkoutRepository:     checkoutRepository,
		userCustomerRepository: userCustomerRepository,
		productRepository:      productRepository,
	}
}

func (cs *checkoutService) CreateCheckout(ctx context.Context, body domain.CheckoutRequest) domain.MessageErr {
	checkout, productCheckouts := body.NewCheckouts()

	ok, err := cs.userCustomerRepository.CheckCustomerExistsByID(ctx, cs.db, checkout.UserCustomerID)
	if err != nil {
		return domain.NewInternalServerError(err.Error())
	}
	if !ok {
		return domain.NewNotFoundError("customerId is not found")
	}

	var productIDs []string
	productQuantities := map[string]int{}
	for _, pc := range productCheckouts {
		productIDs = append(productIDs, pc.ProductID)
		productQuantities[pc.ProductID] = pc.Quantity
	}

	ok, err = cs.productRepository.CheckProductExists(ctx, cs.db, productIDs)
	if err != nil {
		return domain.NewInternalServerError(err.Error())
	}
	if !ok {
		return domain.NewNotFoundError("one of productIds is not found")
	}

	productsStock, err := cs.productRepository.GetProductStockByIDs(ctx, cs.db, productIDs)
	if err != nil {
		return domain.NewInternalServerError(err.Error())
	}
	for _, ps := range productsStock {
		if ps.Stock < productQuantities[ps.ID] {
			return domain.NewBadRequestError(fmt.Sprintf("%s stock is not enough", ps.Name))
		}
	}

	ok, err = cs.productRepository.CheckProductAvailabilities(ctx, cs.db, productIDs)
	if err != nil {
		return domain.NewInternalServerError(err.Error())
	}
	if !ok {
		return domain.NewBadRequestError("one of productIds isAvailable == false")
	}

	totalPrice := 0
	productPrices, err := cs.productRepository.GetProductPriceByIDs(ctx, cs.db, productIDs)
	if err != nil {
		return domain.NewInternalServerError(err.Error())
	}
	for _, ps := range productPrices {
		totalPrice += ps.Price * productQuantities[ps.ID]
	}
	if checkout.Paid < totalPrice {
		return domain.NewBadRequestError(fmt.Sprintf("not enough money, total price is %d", totalPrice))
	}
	change := checkout.Paid - totalPrice
	if *checkout.Change != change {
		return domain.NewBadRequestError(fmt.Sprintf("change is incorrect should be %d", change))
	}

	tx, err := cs.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.NewInternalServerError(err.Error())
	}
	defer tx.Rollback()

	err = cs.checkoutRepository.CreateCheckout(ctx, tx, checkout, productCheckouts)
	if err != nil {
		return domain.NewInternalServerError(err.Error())
	}

	err = cs.checkoutRepository.BulkCreateProductCheckout(ctx, tx, productCheckouts)
	if err != nil {
		return domain.NewInternalServerError(err.Error())
	}

	for _, pc := range productCheckouts {
		err = cs.productRepository.UpdateProductStockByID(ctx, tx, pc.ID, pc.Quantity)
		if err != nil {
			return domain.NewInternalServerError(err.Error())
		}
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
