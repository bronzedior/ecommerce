package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"order/cmd/order/service"
	"order/infrastructure/constant"
	"order/models"
	"time"
)

type OrderUsecase struct {
	OrderService service.OrderService
}

func NewOrderUsecase(orderService service.OrderService) *OrderUsecase {
	return &OrderUsecase{
		OrderService: orderService,
	}
}

func (uc *OrderUsecase) CheckoutOrder(ctx context.Context, param *models.CheckoutRequest) (int64, error) {
	if param.IdempotencyToken != "" {
		isExist, err := uc.OrderService.CheckIdempotency(ctx, param.IdempotencyToken)
		if err != nil {
			return 0, err
		}

		if isExist {
			return 0, errors.New("order already created, please check again")
		}
	}

	if err := uc.validateProducts(param.Items); err != nil {
		return 0, err
	}

	totalQty, totalAmount := uc.calculateOrderSummary(param.Items)
	productJSON, historyJSON, err := uc.constructOrderDetail(param.Items)
	if err != nil {
		return 0, err
	}

	orderDetail := &models.OrderDetail{
		Products:     productJSON,
		OrderHistory: historyJSON,
	}

	order := &models.Order{
		UserID:          param.UserID,
		Amount:          totalAmount,
		TotalQty:        totalQty,
		Status:          constant.OrderStatusCreated,
		PaymentMethod:   param.PaymentMethod,
		ShippingAddress: param.ShippingAddress,
	}

	orderID, err := uc.OrderService.SaveOrderAndOrderDetail(ctx, order, orderDetail)
	if err != nil {
		return 0, err
	}

	if param.IdempotencyToken != "" {
		_ = uc.OrderService.SaveIdempotencyToken(ctx, param.IdempotencyToken)
	}

	return orderID, nil
}

func (uc *OrderUsecase) validateProducts(items []models.CheckoutItem) error {
	seen := map[int64]bool{}
	for _, item := range items {
		if seen[item.ProductID] {
			return fmt.Errorf("duplicate product: %d", item.ProductID)
		}
		seen[item.ProductID] = true

		if item.Quantity <= 0 || item.Quantity > 10000 {
			return fmt.Errorf("invalid quantity for %d", item.ProductID)
		}

		if item.Price <= 0 {
			return fmt.Errorf("invalid price for %d", item.ProductID)
		}
	}
	return nil
}

func (uc *OrderUsecase) calculateOrderSummary(items []models.CheckoutItem) (int, float64) {
	var totalQty int
	var totalAmount float64
	for _, item := range items {
		totalQty += item.Quantity
		totalAmount += float64(item.Quantity) * item.Price
	}
	return totalQty, totalAmount
}

func (uc *OrderUsecase) constructOrderDetail(items []models.CheckoutItem) (string, string, error) {
	productsJSON, _ := json.Marshal(items)
	history := []map[string]interface{}{
		{"status": "created", "timestamp": time.Now()},
	}
	historyJSON, _ := json.Marshal(history)

	return string(productsJSON), string(historyJSON), nil
}

func (uc *OrderUsecase) GetOrderHistory(ctx context.Context, param models.OrderHistoryParam) ([]models.OrderHistoryResponse, error) {
	orderHistory, err := uc.OrderService.GetOrderHistoriesByUserID(ctx, param)
	if err != nil {
		return nil, err
	}

	return orderHistory, nil
}
